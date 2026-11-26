// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatMessage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	transitionTool    = "liveChatMessage-transition"
	transitionConfirm = "Transition live chat message(s) %s to %s"
	transitionShort   = "Transition a live chat message"
	transitionLong    = "Transition a durable live chat event. Use this tool to change the status of a live chat message (e.g., close a poll)."
	transitionExample = `  # Transition a live chat message to closed
  yutu liveChatMessage transition --ids abc123 --status closed`
)

var transitionInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live chat messages to transition",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"status": {
			Type: "string", Enum: []any{"closed"},
			Description: statusUsage,
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: transitionTool, Title: transitionShort, Description: transitionLong,
			InputSchema: transitionInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			transitionTool, cobramcp.ConfirmThen(
				func(input liveChatMessage.LiveChatMessage) string {
					return fmt.Sprintf(
						transitionConfirm, strings.Join(input.Ids, ", "), input.Status,
					)
				},
				func(input liveChatMessage.LiveChatMessage, w io.Writer) error {
					return input.Transition(w)
				},
			),
		),
	)
	liveChatMessageCmd.AddCommand(newTransitionCmd())
}

type transitionOpts struct {
	ids    []string
	status string
	output string
	yes    bool
}

func (o *transitionOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(
		&o.ids, "ids", "i", nil, "IDs of the messages to transition",
	)
	c.Flags().StringVarP(&o.status, "status", "s", "", statusUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
	_ = c.MarkFlagRequired("status")
}

func (o *transitionOpts) build() *liveChatMessage.LiveChatMessage {
	return liveChatMessage.NewLiveChatMessage(
		liveChatMessage.WithIds(o.ids),
		liveChatMessage.WithStatus(o.status),
		liveChatMessage.WithOutput(o.output),
	).(*liveChatMessage.LiveChatMessage)
}

func newTransitionCmd() *cobra.Command {
	var opts transitionOpts

	c := &cobra.Command{
		Use:     "transition",
		Short:   transitionShort,
		Long:    transitionLong,
		Example: transitionExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes,
				fmt.Sprintf(transitionConfirm, strings.Join(opts.ids, ", "), opts.status),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Transition(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

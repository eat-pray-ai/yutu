// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveBroadcast"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	transitionTool    = "liveBroadcast-transition"
	transitionConfirm = "Transition live broadcast(s) %s to %s"
	transitionBsUsage = "Broadcast status to transition to (testing, live, complete)"
	transitionShort   = "Transition a live broadcast"
	transitionLong    = "Transition a live broadcast. Use this tool to change the status of a live broadcast (e.g., go live, end broadcast)."
	transitionExample = `  # Transition a broadcast to live
  yutu liveBroadcast transition --ids broadcast123 --broadcastStatus live
  # End a broadcast
  yutu liveBroadcast transition --ids broadcast123 --broadcastStatus complete`
)

var transitionInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "broadcast_status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live broadcasts to transition",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"broadcast_status": {
			Type: "string", Enum: []any{"testing", "live", "complete"},
			Description: transitionBsUsage,
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","status"]`),
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
				func(input liveBroadcast.LiveBroadcast) string {
					return fmt.Sprintf(
						transitionConfirm, strings.Join(input.Ids, ", "),
						input.BroadcastStatus,
					)
				},
				func(input liveBroadcast.LiveBroadcast, w io.Writer) error {
					return input.Transition(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(newTransitionCmd())
}

type transitionOpts struct {
	ids             []string
	broadcastStatus string
	oboco           string
	obococ          string
	parts           []string
	output          string
	yes             bool
}

func (o *transitionOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(
		&o.ids, "ids", "i", nil, "IDs of the broadcasts to transition",
	)
	c.Flags().StringVarP(
		&o.broadcastStatus, "broadcastStatus", "s", "", transitionBsUsage,
	)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
	_ = c.MarkFlagRequired("broadcastStatus")
}

func (o *transitionOpts) build() *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithIds(o.ids),
		liveBroadcast.WithBroadcastStatus(o.broadcastStatus),
		liveBroadcast.WithOBOCO(o.oboco),
		liveBroadcast.WithOBOCOC(o.obococ),
		liveBroadcast.WithParts(o.parts),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
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
				fmt.Sprintf(
					transitionConfirm, strings.Join(opts.ids, ", "),
					opts.broadcastStatus,
				),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Transition(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

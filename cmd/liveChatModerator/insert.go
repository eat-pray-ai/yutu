// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatModerator"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "liveChatModerator-insert"
	insertConfirm = "Add moderator %s to live chat %s"
	insertShort   = "Add a live chat moderator"
	insertLong    = "Add a live chat moderator. Use this tool to add a moderator to a live chat."
	insertExample = `  # Add a moderator to a live chat
  yutu liveChatModerator insert --liveChatId abc123 --moderatorChannelId UC_xyz`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id", "moderator_channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id":         {Type: "string", Description: lcidUsage},
		"moderator_channel_id": {Type: "string", Description: mcidUsage},
		"parts": {
			Type: "array", Description: "Parts to include in the response",
			Items: &jsonschema.Schema{Type: "string"}, Default: jsontext.Value(`["snippet"]`),
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input liveChatModerator.LiveChatModerator) string {
					return fmt.Sprintf(
						insertConfirm, input.ModeratorChannelId, input.LiveChatId,
					)
				},
				func(input liveChatModerator.LiveChatModerator, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	liveChatModeratorCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	liveChatId         string
	moderatorChannelId string
	parts              []string
	output             string
	yes                bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.liveChatId, "liveChatId", "L", "", lcidUsage)
	c.Flags().StringVarP(
		&o.moderatorChannelId, "moderatorChannelId", "m", "", mcidUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("liveChatId")
	_ = c.MarkFlagRequired("moderatorChannelId")
}

func (o *insertOpts) build() *liveChatModerator.LiveChatModerator {
	return liveChatModerator.NewLiveChatModerator(
		liveChatModerator.WithLiveChatId(o.liveChatId),
		liveChatModerator.WithModeratorChannelId(o.moderatorChannelId),
		liveChatModerator.WithParts(o.parts),
		liveChatModerator.WithOutput(o.output),
	).(*liveChatModerator.LiveChatModerator)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.moderatorChannelId, opts.liveChatId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
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
	deleteTool     = "liveChatMessage-delete"
	deleteConfirm  = "Delete live chat message(s): %s"
	deleteIdsUsage = "IDs of the live chat messages to delete"
	deleteShort    = "Delete live chat messages"
	deleteLong     = "Delete live chat messages. Use this tool to remove messages from a live chat by their IDs."
	deleteExample  = `  # Delete a live chat message by ID
  yutu liveChatMessage delete --ids abc123
  # Delete multiple live chat messages
  yutu liveChatMessage delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(true),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			deleteTool, cobramcp.ConfirmThen(
				func(input liveChatMessage.LiveChatMessage) string {
					return fmt.Sprintf(deleteConfirm, strings.Join(input.Ids, ", "))
				},
				func(input liveChatMessage.LiveChatMessage, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	liveChatMessageCmd.AddCommand(newDeleteCmd())
}

type deleteOpts struct {
	ids []string
	yes bool
}

func (o *deleteOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, deleteIdsUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *deleteOpts) build() *liveChatMessage.LiveChatMessage {
	return liveChatMessage.NewLiveChatMessage(
		liveChatMessage.WithIds(o.ids),
	).(*liveChatMessage.LiveChatMessage)
}

func newDeleteCmd() *cobra.Command {
	var opts deleteOpts

	c := &cobra.Command{
		Use:     "delete",
		Short:   deleteShort,
		Long:    deleteLong,
		Example: deleteExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(deleteConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Delete(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

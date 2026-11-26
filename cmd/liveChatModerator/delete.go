// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"fmt"
	"io"
	"strings"

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
	deleteTool     = "liveChatModerator-delete"
	deleteConfirm  = "Delete live chat moderator(s): %s"
	deleteIdsUsage = "IDs of the live chat moderators to delete"
	deleteShort    = "Delete live chat moderators"
	deleteLong     = "Delete live chat moderators. Use this tool to remove moderators from a live chat by their IDs."
	deleteExample  = `  # Delete a live chat moderator by ID
  yutu liveChatModerator delete --ids abc123
  # Delete multiple live chat moderators
  yutu liveChatModerator delete --ids abc123,def456`
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
				func(input liveChatModerator.LiveChatModerator) string {
					return fmt.Sprintf(deleteConfirm, strings.Join(input.Ids, ", "))
				},
				func(input liveChatModerator.LiveChatModerator, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	liveChatModeratorCmd.AddCommand(newDeleteCmd())
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

func (o *deleteOpts) build() *liveChatModerator.LiveChatModerator {
	return liveChatModerator.NewLiveChatModerator(
		liveChatModerator.WithIds(o.ids),
	).(*liveChatModerator.LiveChatModerator)
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

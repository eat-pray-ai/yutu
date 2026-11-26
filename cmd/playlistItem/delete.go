// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "playlistItem-delete"
	deleteIdsUsage = "IDs of the playlist items to delete"
	deleteConfirm  = "Delete playlist item(s): %s"
	deleteShort    = "Delete playlist items"
	deleteLong     = "Delete playlist items. Use this tool to delete playlist items by IDs."
	deleteExample  = `  # Delete a playlist item by ID
  yutu playlistItem delete --ids abc123
  # Delete multiple playlist items
  yutu playlistItem delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
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
				func(input playlistItem.PlaylistItem) string {
					return fmt.Sprintf(deleteConfirm, strings.Join(input.Ids, ", "))
				},
				func(input playlistItem.PlaylistItem, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	playlistItemCmd.AddCommand(newDeleteCmd())
}

type deleteOpts struct {
	ids                    []string
	onBehalfOfContentOwner string
	yes                    bool
}

func (o *deleteOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, deleteIdsUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *deleteOpts) build() *playlistItem.PlaylistItem {
	return playlistItem.NewPlaylistItem(
		playlistItem.WithIds(o.ids),
		playlistItem.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
	).(*playlistItem.PlaylistItem)
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

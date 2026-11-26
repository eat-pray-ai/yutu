// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool    = "playlistImage-delete"
	deleteConfirm = "Delete playlist image(s): %s"
	deleteShort   = "Delete playlist images"
	deleteLong    = "Delete playlist images. Use this tool to delete playlist images by IDs."
	deleteExample = `  # Delete a playlist image by ID
  yutu playlistImage delete --ids abc123
  # Delete multiple playlist images
  yutu playlistImage delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
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
				func(input playlistImage.PlaylistImage) string {
					return fmt.Sprintf(deleteConfirm, strings.Join(input.Ids, ", "))
				},
				func(input playlistImage.PlaylistImage, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	playlistImageCmd.AddCommand(newDeleteCmd())
}

type deleteOpts struct {
	ids                    []string
	onBehalfOfContentOwner string
	yes                    bool
}

func (o *deleteOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, idsUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *deleteOpts) build() *playlistImage.PlaylistImage {
	return playlistImage.NewPlaylistImage(
		playlistImage.WithIds(o.ids),
		playlistImage.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
	).(*playlistImage.PlaylistImage)
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

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
	updateTool    = "playlistImage-update"
	updateConfirm = "Update playlist image: %s"
	updateShort   = "Update a playlist image"
	updateLong    = "Update a playlist image. Use this tool to update a playlist image."
	updateExample = `  # Update a playlist image
  yutu playlistImage update --playlistId PLxxx --type hero --width 2048 --height 1152`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"playlist_id"},
	Properties: map[string]*jsonschema.Schema{
		"playlist_id": {Type: "string", Description: pidUsage},
		"type":        {Type: "string", Description: typeUsage},
		"height": {
			Type: "number", Description: heightUsage,
			Minimum: new(float64(0)),
		},
		"width": {
			Type: "number", Description: widthUsage,
			Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input playlistImage.PlaylistImage) string {
					return fmt.Sprintf(updateConfirm, input.PlaylistId)
				},
				func(input playlistImage.PlaylistImage, w io.Writer) error {
					playlistImage.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	playlistImageCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	playlistId string
	type_      string
	height     int64
	width      int64
	output     string
	oboco      string
	obococ     string
	yes        bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.playlistId, "playlistId", "P", "", pidUsage)
	c.Flags().StringVarP(&o.type_, "type", "t", "", typeUsage)
	c.Flags().Int64VarP(&o.height, "height", "H", 0, heightUsage)
	c.Flags().Int64VarP(&o.width, "width", "W", 0, widthUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("playlistId")
}

func (o *updateOpts) build() *playlistImage.PlaylistImage {
	return playlistImage.NewPlaylistImage(
		playlistImage.WithPlaylistId(o.playlistId),
		playlistImage.WithType(o.type_),
		playlistImage.WithHeight(o.height),
		playlistImage.WithWidth(o.width),
		playlistImage.WithMaxResults(1),
		playlistImage.WithOutput(o.output),
		playlistImage.WithOBOCO(o.oboco),
		playlistImage.WithOBOCOC(o.obococ),
	).(*playlistImage.PlaylistImage)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, opts.playlistId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

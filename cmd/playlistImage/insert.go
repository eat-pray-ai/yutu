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
	insertTool    = "playlistImage-insert"
	insertConfirm = "Upload playlist image %s for playlist %s"
	insertShort   = "Upload a playlist image"
	insertLong    = "Upload a playlist image. Use this tool to upload a playlist image."
	insertExample = `  # Insert a playlist cover image
  yutu playlistImage insert --file cover.jpg --playlistId PLxxx
  # Insert a hero image
  yutu playlistImage insert --file cover.png --playlistId PLxxx --type hero
  # Insert an image with custom dimensions
  yutu playlistImage insert --file cover.jpg --playlistId PLxxx --width 2048 --height 1152`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "playlist_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":        {Type: "string", Description: fileUsage},
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input playlistImage.PlaylistImage) string {
					return fmt.Sprintf(insertConfirm, input.File, input.PlaylistId)
				},
				func(input playlistImage.PlaylistImage, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	playlistImageCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	file       string
	playlistId string
	type_      string
	height     int64
	width      int64
	output     string
	oboco      string
	obococ     string
	yes        bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(&o.playlistId, "playlistId", "p", "", pidUsage)
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
	_ = c.MarkFlagRequired("file")
	_ = c.MarkFlagRequired("playlistId")
}

func (o *insertOpts) build() *playlistImage.PlaylistImage {
	return playlistImage.NewPlaylistImage(
		playlistImage.WithFile(o.file),
		playlistImage.WithPlaylistId(o.playlistId),
		playlistImage.WithType(o.type_),
		playlistImage.WithHeight(o.height),
		playlistImage.WithWidth(o.width),
		playlistImage.WithOutput(o.output),
		playlistImage.WithOBOCO(o.oboco),
		playlistImage.WithOBOCOC(o.obococ),
	).(*playlistImage.PlaylistImage)
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
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.file, opts.playlistId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

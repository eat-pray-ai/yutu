// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	setTool    = "thumbnail-set"
	setConfirm = "Set thumbnail %s for video %s"
	setShort   = "Set a thumbnail for a video"
	setLong    = "Set a thumbnail for a video. Use this tool to set a thumbnail for a video."
	setExample = `  # Set a thumbnail for a video
  yutu thumbnail set --file image.jpg --videoId dQw4w9WgXcQ
  # Set a thumbnail with JSON output
  yutu thumbnail set --file image.png --videoId dQw4w9WgXcQ --output json`
)

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":     {Type: "string", Description: fileUsage},
		"video_id": {Type: "string", Description: vidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: setTool, Title: setShort, Description: setLong,
			InputSchema: setInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			setTool, cobramcp.ConfirmThen(
				func(input thumbnail.Thumbnail) string {
					return fmt.Sprintf(setConfirm, input.File, input.VideoId)
				},
				func(input thumbnail.Thumbnail, w io.Writer) error {
					return input.Set(w)
				},
			),
		),
	)
	thumbnailCmd.AddCommand(newSetCmd())
}

type setOpts struct {
	file    string
	videoId string
	output  string
	yes     bool
}

func (o *setOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("file")
	_ = c.MarkFlagRequired("videoId")
}

func (o *setOpts) build() *thumbnail.Thumbnail {
	return thumbnail.NewThumbnail(
		thumbnail.WithFile(o.file),
		thumbnail.WithVideoId(o.videoId),
		thumbnail.WithOutput(o.output),
	).(*thumbnail.Thumbnail)
}

func newSetCmd() *cobra.Command {
	var opts setOpts

	c := &cobra.Command{
		Use:     "set",
		Short:   setShort,
		Long:    setLong,
		Example: setExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(setConfirm, opts.file, opts.videoId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Set(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

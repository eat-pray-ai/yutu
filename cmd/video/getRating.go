// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	getRatingTool    = "video-getRating"
	grIdsUsage       = "IDs of the videos to get the rating for"
	getRatingShort   = "Get video ratings"
	getRatingLong    = "Get video ratings. Use this tool to get video ratings by IDs."
	getRatingExample = `  # Get rating of a video
  yutu video getRating --ids dQw4w9WgXcQ
  # Get ratings of multiple videos in JSON format
  yutu video getRating --ids dQw4w9WgXcQ,abc123 --output json`
)

var getRatingInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: grIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: getRatingTool, Title: getRatingShort, Description: getRatingLong,
			InputSchema: getRatingInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cobramcp.GenToolHandler(
			getRatingTool, func(input video.Video, writer io.Writer) error {
				return input.GetRating(writer)
			},
		),
	)
	videoCmd.AddCommand(newGetRatingCmd())
}

type getRatingOpts struct {
	ids    []string
	oboco  string
	output string
}

func (o *getRatingOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, grIdsUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *getRatingOpts) build() *video.Video {
	return video.NewVideo(
		video.WithIds(o.ids),
		video.WithOBOCO(o.oboco),
		video.WithOutput(o.output),
		video.WithService(nil),
	).(*video.Video)
}

func newGetRatingCmd() *cobra.Command {
	var opts getRatingOpts

	c := &cobra.Command{
		Use:     "getRating",
		Short:   getRatingShort,
		Long:    getRatingLong,
		Example: getRatingExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().GetRating(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

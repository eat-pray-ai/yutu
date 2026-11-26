// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	rateTool     = "video-rate"
	rateIdsUsage = "IDs of the videos to rate"
	rateRUsage   = "like|dislike|none"
	rateConfirm  = "Rate video(s) %s as %s"
	rateShort    = "Rate a video"
	rateLong     = "Rate a video. Use this tool to rate a video."
	rateExample  = `  # Like a video
  yutu video rate --ids dQw4w9WgXcQ --rating like
  # Dislike multiple videos
  yutu video rate --ids dQw4w9WgXcQ,abc123 --rating dislike
  # Remove rating from a video
  yutu video rate --ids dQw4w9WgXcQ --rating none`
)

var rateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "rating"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: rateIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"rating": {
			Type: "string", Description: rateRUsage,
			Enum: []any{"like", "dislike", "none"},
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: rateTool, Title: rateShort, Description: rateLong,
			InputSchema: rateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			rateTool, cobramcp.ConfirmThen(
				func(input video.Video) string {
					return fmt.Sprintf(
						rateConfirm, strings.Join(input.Ids, ", "), input.Rating,
					)
				},
				func(input video.Video, w io.Writer) error {
					return input.Rate(w)
				},
			),
		),
	)
	videoCmd.AddCommand(newRateCmd())
}

type rateOpts struct {
	ids    []string
	rating string
	yes    bool
}

func (o *rateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, rateIdsUsage)
	c.Flags().StringVarP(&o.rating, "rating", "r", "", rateRUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
	_ = c.MarkFlagRequired("rating")
}

func (o *rateOpts) build() *video.Video {
	return video.NewVideo(
		video.WithIds(o.ids),
		video.WithRating(o.rating),
	).(*video.Video)
}

func newRateCmd() *cobra.Command {
	var opts rateOpts

	c := &cobra.Command{
		Use:     "rate",
		Short:   rateShort,
		Long:    rateLong,
		Example: rateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(rateConfirm, strings.Join(opts.ids, ", "), opts.rating),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Rate(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

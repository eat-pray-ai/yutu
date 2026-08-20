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
	rateShort    = "Rate a video"
	rateLong     = "Rate a video. Use this tool to rate a video."
	rateExample  = `# Like a video
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
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
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
		}, cobramcp.GenToolHandler(
			rateTool, func(input video.Video, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Rate(writer)
			},
		),
	)
	videoCmd.AddCommand(rateCmd)

	rateCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, rateIdsUsage)
	rateCmd.Flags().StringVarP(&rating, "rating", "r", "", rateRUsage)
	rateCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = rateCmd.MarkFlagRequired("ids")
	_ = rateCmd.MarkFlagRequired("rating")
}

var rateCmd = &cobra.Command{
	Use:     "rate",
	Short:   rateShort,
	Long:    rateLong,
	Example: rateExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf(
			"Would rate video(s): %s as %s", strings.Join(ids, ", "), rating,
		)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithRating(rating),
		)
		utils.HandleCmdError(input.Rate(c.OutOrStdout()), c)
	},
}

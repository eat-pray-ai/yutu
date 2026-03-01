// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	rateTool     = "video-rate"
	rateShort    = "Rate a video on YouTube"
	rateLong     = "Rate a video on YouTube, with the specified rating\n\nExamples:\n  yutu video rate --ids dQw4w9WgXcQ --rating like\n  yutu video rate --ids dQw4w9WgXcQ,abc123 --rating dislike\n  yutu video rate --ids dQw4w9WgXcQ --rating none"
	rateIdsUsage = "IDs of the videos to rate"
	rateRUsage   = "like|dislike|none"
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
			Enum: []any{"like", "dislike", "none", ""},
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
		}, cobramcp.GenToolHandler(
			rateTool, func(input video.Video, writer io.Writer) error {
				return input.Rate(writer)
			},
		),
	)
	videoCmd.AddCommand(rateCmd)

	rateCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, rateIdsUsage)
	rateCmd.Flags().StringVarP(&rating, "rating", "r", "", rateRUsage)

	_ = rateCmd.MarkFlagRequired("ids")
	_ = rateCmd.MarkFlagRequired("rating")
}

var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: rateShort,
	Long:  rateLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithRating(rating),
		)
		utils.HandleCmdError(input.Rate(cmd.OutOrStdout()), cmd)
	},
}

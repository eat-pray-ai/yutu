// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type rateIn struct {
	Ids    []string `json:"ids"`
	Rating string   `json:"rating"`
}

var rateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "rating"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: rateIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"rating": {
			Type: "string", Enum: []any{"like", "dislike", "none", ""},
			Description: rateRUsage, Default: json.RawMessage(`""`),
		},
	},
}

const (
	rateShort    = "Rate a video on YouTube"
	rateLong     = "Rate a video on YouTube, with the specified rating"
	rateIdsUsage = "IDs of the videos to rate"
	rateRUsage   = "like, dislike, or none"
)

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "video-rate", Title: rateShort, Description: rateLong,
			InputSchema: rateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, rateHandler,
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
		err := rate(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func rateHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input rateIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	rating = input.Rating

	slog.InfoContext(ctx, "video rate started")

	var writer bytes.Buffer
	err := rate(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video rate failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "video rate completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func rate(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithRating(rating),
		video.WithService(nil),
	)

	return v.Rate(writer)
}

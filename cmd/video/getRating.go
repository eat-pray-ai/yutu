// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	getRatingTool  = "video-getRating"
	getRatingShort = "Get the rating of videos"
	getRatingLong  = "Get the rating of videos by ids"
	grIdsUsage     = "IDs of the videos to get the rating for"
)

var getRatingInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: grIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of_content_owner": {Type: "string"},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: getRatingTool, Title: getRatingShort, Description: getRatingLong,
			InputSchema: getRatingInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, getRatingHandler,
	)
	videoCmd.AddCommand(getRatingCmd)

	getRatingCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, grIdsUsage)
	getRatingCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	getRatingCmd.Flags().StringVarP(&output, "output", "o", "", pkg.TableUsage)
	getRatingCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = getRatingCmd.MarkFlagRequired("ids")
}

var getRatingCmd = &cobra.Command{
	Use:   "getRating",
	Short: getRatingShort,
	Long:  getRatingLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithOutput(output),
			video.WithJsonpath(jsonpath),
			video.WithService(nil),
		)
		err := input.GetRating(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func getRatingHandler(
	ctx context.Context, req *mcp.CallToolRequest, input video.Video,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: getRatingTool, MinInterval: time.Second,
			},
		),
	)

	var writer bytes.Buffer
	err := input.GetRating(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "comment-insert"
	insertShort    = "Insert a comment"
	insertLong     = "Insert a comment to a video"
	insertPidUsage = "ID of the parent comment"
)

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"author_channel_id", "channel_id", "parent_id", "text_original", "video_id",
	},
	Properties: map[string]*jsonschema.Schema{
		"author_channel_id": {Type: "string", Description: acidUsage},
		"channel_id":        {Type: "string", Description: cidUsage},
		"can_rate": {
			Type:        "boolean",
			Description: crUsage, Default: json.RawMessage(`false`),
		},
		"parent_id":     {Type: "string", Description: insertPidUsage},
		"text_original": {Type: "string", Description: toUsage},
		"video_id":      {Type: "string", Description: vidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	commentCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", cidUsage,
	)
	insertCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	insertCmd.Flags().StringVarP(
		&parentId, "parentId", "P", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("parentId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := comment.NewComment(
			comment.WithAuthorChannelId(authorChannelId),
			comment.WithChannelId(channelId),
			comment.WithCanRate(canRate),
			comment.WithParentId(parentId),
			comment.WithTextOriginal(textOriginal),
			comment.WithVideoId(videoId),
			comment.WithOutput(output),
			comment.WithJsonpath(jsonpath),
		)
		err := input.Insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input comment.Comment,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	var writer bytes.Buffer
	err := input.Insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

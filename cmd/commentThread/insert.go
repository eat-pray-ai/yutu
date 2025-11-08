// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "commentThread-insert"
	insertShort    = "Insert a new comment thread"
	insertLong     = "Insert a new comment thread"
	insertVidUsage = "ID of the video"
)

type insertIn struct {
	AuthorChannelId string `json:"authorChannelId"`
	ChannelId       string `json:"channelId"`
	TextOriginal    string `json:"textOriginal"`
	VideoId         string `json:"videoId"`
	Output          string `json:"output"`
	Jsonpath        string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"authorChannelId", "channelId", "textOriginal", "videoId",
		"output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"authorChannelId": {
			Type: "string", Description: acidUsage,
			Default: json.RawMessage(`""`),
		},
		"channelId": {
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"textOriginal": {
			Type: "string", Description: toUsage,
			Default: json.RawMessage(`""`),
		},
		"videoId": {
			Type: "string", Description: insertVidUsage,
			Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
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
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", toUsage)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", insertVidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	authorChannelId = input.AuthorChannelId
	channelId = input.ChannelId
	textOriginal = input.TextOriginal
	videoId = input.VideoId
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	ct := commentThread.NewCommentThread(
		commentThread.WithAuthorChannelId(authorChannelId),
		commentThread.WithChannelId(channelId),
		commentThread.WithTextOriginal(textOriginal),
		commentThread.WithVideoId(videoId),
		commentThread.WithService(nil),
	)

	return ct.Insert(output, jpath, writer)
}

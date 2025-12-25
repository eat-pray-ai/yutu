// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
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

type insertIn struct {
	AuthorChannelId string `json:"authorChannelId"`
	ChannelId       string `json:"channelId"`
	CanRate         string `json:"canRate"`
	ParentId        string `json:"parentId"`
	TextOriginal    string `json:"textOriginal"`
	VideoId         string `json:"videoId"`
	Output          string `json:"output"`
	Jsonpath        string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"authorChannelId", "channelId", "parentId", "textOriginal", "videoId",
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
		"canRate": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: crUsage, Default: json.RawMessage(`""`),
		},
		"parentId": {
			Type: "string", Description: insertPidUsage,
			Default: json.RawMessage(`""`),
		},
		"textOriginal": {
			Type: "string", Description: toUsage,
			Default: json.RawMessage(`""`),
		},
		"videoId": {
			Type: "string", Description: vidUsage,
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
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

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
	canRate = utils.BoolPtr(input.CanRate)
	parentId = input.ParentId
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
	c := comment.NewComment(
		comment.WithAuthorChannelId(authorChannelId),
		comment.WithChannelId(channelId),
		comment.WithCanRate(canRate),
		comment.WithParentId(parentId),
		comment.WithTextOriginal(textOriginal),
		comment.WithVideoId(videoId),
		comment.WithService(nil),
	)

	return c.Insert(output, jpath, writer)
}

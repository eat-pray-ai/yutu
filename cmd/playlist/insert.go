// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "playlist-insert"
	insertShort    = "Create a new playlist"
	insertLong     = "Create a new playlist, with the specified title, description, tags, etc"
	insertCidUsage = "Channel id of the playlist"
)

type insertIn struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Language    string   `json:"language"`
	ChannelId   string   `json:"channelId"`
	Privacy     string   `json:"privacy"`
	Output      string   `json:"output"`
	Jsonpath    string   `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"title", "channelId", "privacy"},
	Properties: map[string]*jsonschema.Schema{
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"tags": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: tagsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"language": {
			Type: "string", Description: languageUsage,
			Default: json.RawMessage(`""`),
		},
		"channelId": {
			Type: "string", Description: insertCidUsage,
			Default: json.RawMessage(`""`),
		},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage, Default: json.RawMessage(`""`),
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
	playlistCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("title")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("privacy")
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

	title = input.Title
	description = input.Description
	tags = input.Tags
	language = input.Language
	channelId = input.ChannelId
	privacy = input.Privacy
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithTitle(title),
		playlist.WithDescription(description),
		playlist.WithTags(tags),
		playlist.WithLanguage(language),
		playlist.WithChannelId(channelId),
		playlist.WithPrivacy(privacy),
		playlist.WithService(nil),
	)

	return p.Insert(output, jsonpath, writer)
}

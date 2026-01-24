// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool  = "caption-insert"
	insertShort = "Insert caption"
	insertLong  = "Insert caption to a video"
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "videoId"},
	Properties: map[string]*jsonschema.Schema{
		"file": {Type: "string", Description: fileUsage},
		"audio_track_type": {
			Type: "string", Description: attUsage,
			Enum:    []any{"unknown", "primary", "commentary", "descriptive"},
			Default: json.RawMessage(`"unknown"`),
		},
		"is_auto_synced": {
			Type: "boolean", Description: iasUsage,
			Enum: []any{true, false},
		},
		"is_cc": {
			Type: "boolean", Description: iscUsage,
			Enum: []any{true, false},
		},
		"is_draft": {
			Type: "boolean", Description: isdUsage,
			Enum: []any{true, false},
		},
		"is_easy_reader": {
			Type: "boolean", Description: iserUsage,
			Enum: []any{true, false},
		},
		"is_large": {
			Type: "boolean", Description: islUsage,
			Enum: []any{true, false},
		},
		"language": {Type: "string", Description: languageUsage},
		"name":     {Type: "string", Description: nameUsage},
		"track_kind": {
			Type: "string", Description: tkUsage,
			Enum:    []any{"standard", "ASR", "forced"},
			Default: json.RawMessage(`"standard"`),
		},
		"video_id": {
			Type: "string", Description: vidUsage,
			Default: json.RawMessage(`""`),
		},
		"on_behalf_of":               {Type: "string"},
		"on_behalf_of_content_owner": {Type: "string"},
		"output": {
			Type: "string", Description: pkg.SilentUsage,
			Enum:    []any{"json", "yaml", "silent", ""},
			Default: json.RawMessage(`"yaml"`),
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
	captionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	insertCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	insertCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	insertCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	insertCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	insertCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	insertCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := caption.NewCation(
			caption.WithFile(file),
			caption.WithAudioTrackType(audioTrackType),
			caption.WithIsAutoSynced(isAutoSynced),
			caption.WithIsCC(isCC),
			caption.WithIsDraft(isDraft),
			caption.WithIsEasyReader(isEasyReader),
			caption.WithIsLarge(isLarge),
			caption.WithLanguage(language),
			caption.WithName(name),
			caption.WithTrackKind(trackKind),
			caption.WithVideoId(videoId),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithOutput(output),
			caption.WithJsonpath(jsonpath),
		)
		if err := input.Insert(cmd.OutOrStdout()); err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input caption.Caption,
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
	if err := input.Insert(&writer); err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

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
	smsTool  = "comment-setModerationStatus"
	smsShort = "Set YouTube comments moderation status"
	smsLong  = "Set YouTube comments moderation status by ids"
)

var setModerationStatusInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "moderation_status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"moderation_status": {
			Type: "string", Description: msUsage,
			Enum: []any{"heldForReview", "published", "rejected", ""},
		},
		"ban_author": {
			Type: "boolean", Description: baUsage, Default: json.RawMessage(`false`),
		},
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
			Name: smsTool, Title: smsShort, Description: smsLong,
			InputSchema: setModerationStatusInSchema,
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, setModerationStatusHandler,
	)
	commentCmd.AddCommand(setModerationStatusCmd)

	setModerationStatusCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, idsUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "s", "", msUsage,
	)
	setModerationStatusCmd.Flags().BoolVarP(
		banAuthor, "banAuthor", "A", false, baUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&output, "output", "o", "", pkg.SilentUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&jsonpath, "jsonpath", "j", "", pkg.JPUsage,
	)

	_ = setModerationStatusCmd.MarkFlagRequired("ids")
	_ = setModerationStatusCmd.MarkFlagRequired("moderationStatus")
}

var setModerationStatusCmd = &cobra.Command{
	Use:   "setModerationStatus",
	Short: smsShort,
	Long:  smsLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := comment.NewComment(
			comment.WithIds(ids),
			comment.WithModerationStatus(moderationStatus),
			comment.WithBanAuthor(banAuthor),
			comment.WithOutput(output),
			comment.WithJsonpath(jsonpath),
		)
		err := input.SetModerationStatus(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setModerationStatusHandler(
	ctx context.Context, req *mcp.CallToolRequest, input comment.Comment,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: smsTool, MinInterval: time.Second},
		),
	)

	var writer bytes.Buffer
	err := input.SetModerationStatus(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

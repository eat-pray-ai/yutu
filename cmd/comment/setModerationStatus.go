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
	smsTool  = "comment-setModerationStatus"
	smsShort = "Set YouTube comments moderation status"
	smsLong  = "Set YouTube comments moderation status by ids"
)

type setModerationStatusIn struct {
	IDs              []string `json:"ids"`
	ModerationStatus string   `json:"moderationStatus"`
	BanAuthor        string   `json:"banAuthor"`
	Output           string   `json:"output"`
	Jsonpath         string   `json:"jsonpath"`
}

var setModerationStatusInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "moderationStatus", "banAuthor",
		"output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: idsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"moderationStatus": {
			Type: "string", Enum: []any{"heldForReview", "published", "rejected", ""},
			Description: msUsage, Default: json.RawMessage(`""`),
		},
		"banAuthor": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: baUsage, Default: json.RawMessage(`""`),
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
		&jpath, "jsonpath", "j", "", pkg.JPUsage,
	)

	_ = setModerationStatusCmd.MarkFlagRequired("ids")
	_ = setModerationStatusCmd.MarkFlagRequired("moderationStatus")
}

var setModerationStatusCmd = &cobra.Command{
	Use:   "setModerationStatus",
	Short: smsShort,
	Long:  smsLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := setModerationStatus(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setModerationStatusHandler(
	ctx context.Context, req *mcp.CallToolRequest, input setModerationStatusIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: smsTool, MinInterval: time.Second},
		),
	)

	ids = input.IDs
	moderationStatus = input.ModerationStatus
	banAuthor = utils.BoolPtr(input.BanAuthor)
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := setModerationStatus(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func setModerationStatus(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithModerationStatus(moderationStatus),
		comment.WithBanAuthor(banAuthor),
	)

	return c.SetModerationStatus(output, jpath, writer)
}

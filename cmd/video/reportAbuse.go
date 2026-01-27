// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	reportAbuseTool  = "video-reportAbuse"
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video"
	raIdsUsage       = "IDs of the videos to report abuse on"
	raLangUsage      = "Language that the content was viewed in"
)

var reportAbuseInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "reason_id"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: raIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"reason_id":                  {Type: "string", Description: ridUsage},
		"secondary_reason_id":        {Type: "string", Description: sridUsage},
		"comments":                   {Type: "string", Description: commentsUsage},
		"language":                   {Type: "string", Description: raLangUsage},
		"on_behalf_of_content_owner": {Type: "string"},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: reportAbuseTool, Title: reportAbuseShort,
			Description: reportAbuseLong,
			InputSchema: reportAbuseInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, reportAbuseHandler,
	)
	videoCmd.AddCommand(reportAbuseCmd)

	reportAbuseCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, raIdsUsage,
	)
	reportAbuseCmd.Flags().StringVarP(&reasonId, "reasonId", "r", "", ridUsage)
	reportAbuseCmd.Flags().StringVarP(
		&secondaryReasonId, "secondaryReasonId", "s", "", sridUsage,
	)
	reportAbuseCmd.Flags().StringVarP(
		&comments, "comments", "c", "", commentsUsage,
	)
	reportAbuseCmd.Flags().StringVarP(&language, "language", "l", "", raLangUsage)
	reportAbuseCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = reportAbuseCmd.MarkFlagRequired("ids")
	_ = reportAbuseCmd.MarkFlagRequired("reasonId")
}

var reportAbuseCmd = &cobra.Command{
	Use:   "reportAbuse",
	Short: reportAbuseShort,
	Long:  reportAbuseLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithReasonId(reasonId),
			video.WithSecondaryReasonId(secondaryReasonId),
			video.WithComments(comments),
			video.WithLanguage(language),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		err := input.ReportAbuse(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func reportAbuseHandler(
	ctx context.Context, req *mcp.CallToolRequest, input video.Video,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: reportAbuseTool, MinInterval: time.Second,
			},
		),
	)

	var writer bytes.Buffer
	err := input.ReportAbuse(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

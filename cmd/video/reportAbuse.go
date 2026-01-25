// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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

type reportAbuseIn struct {
	Ids                    []string `json:"ids"`
	ReasonId               string   `json:"reasonId"`
	SecondaryReasonId      string   `json:"secondaryReasonId"`
	Comments               string   `json:"comments"`
	Language               string   `json:"language"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
}

var reportAbuseInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "reasonId"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: raIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"reasonId":               {Type: "string", Description: ridUsage},
		"secondaryReasonId":      {Type: "string", Description: sridUsage},
		"comments":               {Type: "string", Description: commentsUsage},
		"language":               {Type: "string", Description: raLangUsage},
		"onBehalfOfContentOwner": {Type: "string"},
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
		err := reportAbuse(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func reportAbuseHandler(
	ctx context.Context, req *mcp.CallToolRequest, input reportAbuseIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: reportAbuseTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.Ids
	reasonId = input.ReasonId
	secondaryReasonId = input.SecondaryReasonId
	comments = input.Comments
	language = input.Language
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	var writer bytes.Buffer
	err := reportAbuse(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func reportAbuse(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIds(ids),
		video.WithReasonId(reasonId),
		video.WithSecondaryReasonId(secondaryReasonId),
		video.WithComments(comments),
		video.WithLanguage(language),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithService(nil),
	)

	return v.ReportAbuse(writer)
}

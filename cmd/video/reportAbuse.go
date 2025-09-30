// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
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
	Type: "object",
	Required: []string{
		"ids", "reasonId", "secondaryReasonId", "comments",
		"language", "onBehalfOfContentOwner",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: raIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"reasonId": {
			Type: "string", Description: ridUsage,
			Default: json.RawMessage(`""`),
		},
		"secondaryReasonId": {
			Type: "string", Description: sridUsage,
			Default: json.RawMessage(`""`),
		},
		"comments": {
			Type: "string", Description: commentsUsage,
			Default: json.RawMessage(`""`),
		},
		"language": {
			Type: "string", Description: raLangUsage,
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
	},
}

const (
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video"
	raIdsUsage       = "IDs of the videos to report abuse on"
	raLangUsage      = "Language that the content was viewed in"
)

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "video-reportAbuse", Title: reportAbuseShort,
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
	ctx context.Context, _ *mcp.CallToolRequest, input reportAbuseIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	reasonId = input.ReasonId
	secondaryReasonId = input.SecondaryReasonId
	comments = input.Comments
	language = input.Language
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	slog.InfoContext(ctx, "video reportAbuse started")

	var writer bytes.Buffer
	err := reportAbuse(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video reportAbuse failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "video reportAbuse completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func reportAbuse(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithReasonId(reasonId),
		video.WithSecondaryReasonId(secondaryReasonId),
		video.WithComments(comments),
		video.WithLanguage(language),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithService(nil),
	)

	return v.ReportAbuse(writer)
}

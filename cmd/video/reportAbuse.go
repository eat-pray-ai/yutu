// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	reportAbuseTool  = "video-reportAbuse"
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video\n\nExamples:\n  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V\n  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId V --secondaryReasonId V1 --language en\n  yutu video reportAbuse --ids dQw4w9WgXcQ --reasonId N --comments 'Spam content'"
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
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: reportAbuseTool, Title: reportAbuseShort,
			Description: reportAbuseLong,
			InputSchema: reportAbuseInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			reportAbuseTool, func(input video.Video, writer io.Writer) error {
				return input.ReportAbuse(writer)
			},
		),
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
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
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
		utils.HandleCmdError(input.ReportAbuse(cmd.OutOrStdout()), cmd)
	},
}

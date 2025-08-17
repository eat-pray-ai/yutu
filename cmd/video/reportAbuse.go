package video

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video"
	raIdsUsage       = "IDs of the videos to report abuse on"
	raLangUsage      = "Language that the content was viewed in"
)

func init() {
	cmd.MCP.AddTool(reportAbuseTool, reportAbuseHandler)
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

var reportAbuseTool = mcp.NewTool(
	"video-reportAbuse",
	mcp.WithTitleAnnotation(reportAbuseShort),
	mcp.WithDescription(reportAbuseLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(raIdsUsage), mcp.Required(),
	),
	mcp.WithString(
		"reasonId", mcp.DefaultString(""), mcp.Description(ridUsage), mcp.Required(),
	),
	mcp.WithString(
		"secondaryReasonId", mcp.DefaultString(""), mcp.Description(sridUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"comments", mcp.DefaultString(""), mcp.Description(commentsUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"language", mcp.DefaultString(""), mcp.Description(raLangUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""), mcp.Description(""),
		mcp.Required(),
	),
)

func reportAbuseHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	reasonId, _ = args["reasonId"].(string)
	secondaryReasonId, _ = args["secondaryReasonId"].(string)
	comments, _ = args["comments"].(string)
	language, _ = args["language"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)

	slog.InfoContext(ctx, "video reportAbuse started")

	var writer bytes.Buffer
	err := reportAbuse(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video reportAbuse failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "video reportAbuse completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
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

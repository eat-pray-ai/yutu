package comment

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	smsShort = "Set YouTube comments moderation status"
	smsLong  = "Set YouTube comments moderation status by ids"
)

func init() {
	cmd.MCP.AddTool(setModerationStatusTool, setModerationStatusHandler)
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

var setModerationStatusTool = mcp.NewTool(
	"comment-setModerationStatus",
	mcp.WithTitleAnnotation(smsShort),
	mcp.WithDescription(smsLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(idsUsage), mcp.Required(),
	),
	mcp.WithString(
		"moderationStatus", mcp.Enum("heldForReview", "published", "rejected"),
		mcp.DefaultString(""), mcp.Description(msUsage), mcp.Required(),
	),
	mcp.WithString(
		"banAuthor", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(baUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "silent", ""),
		mcp.DefaultString("yaml"), mcp.Description(pkg.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(pkg.JPUsage), mcp.Required(),
	),
)

func setModerationStatusHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	moderationStatus, _ = args["moderationStatus"].(string)
	banAuthorRaw, _ := args["banAuthor"].(string)
	banAuthor = utils.BoolPtr(banAuthorRaw)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "comment setModerationStatus started")

	var writer bytes.Buffer
	err := setModerationStatus(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "comment setModerationStatus failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "comment setModerationStatus completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func setModerationStatus(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithModerationStatus(moderationStatus),
		comment.WithBanAuthor(banAuthor),
	)

	return c.SetModerationStatus(output, jpath, writer)
}

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
	updateShort   = "Update a comment"
	updateLong    = "Update a comment on a video"
	updateIdUsage = "ID of the comment"
)

func init() {
	cmd.MCP.AddTool(updateTool, updateHandler)
	commentCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	updateCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	updateCmd.Flags().StringVarP(
		&viewerRating, "viewerRating", "r", "", vrUsage,
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var updateTool = mcp.NewTool(
	"comment-update",
	mcp.WithTitleAnnotation(updateShort),
	mcp.WithDescription(updateLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(updateIdUsage), mcp.Required(),
	),
	mcp.WithString(
		"canRate", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(crUsage), mcp.Required(),
	),
	mcp.WithString(
		"textOriginal", mcp.DefaultString(""),
		mcp.Description(toUsage), mcp.Required(),
	),
	mcp.WithString(
		"viewerRating", mcp.Enum("none", "like", "dislike"),
		mcp.DefaultString(""), mcp.Description(vrUsage), mcp.Required(),
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

func updateHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	canRateRaw, _ := args["canRate"].(string)
	canRate = utils.BoolPtr(canRateRaw)
	textOriginal, _ = args["textOriginal"].(string)
	viewerRating, _ = args["viewerRating"].(string)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "comment update started")

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "comment update failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "comment update completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func update(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithCanRate(canRate),
		comment.WithTextOriginal(textOriginal),
		comment.WithViewerRating(viewerRating),
		comment.WithMaxResults(1),
		comment.WithService(nil),
	)

	return c.Update(output, jpath, writer)
}

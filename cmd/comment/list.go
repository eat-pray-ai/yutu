package comment

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List YouTube comments"
	listLong     = "List YouTube comments by ids"
	listPidUsage = "Returns replies to the specified comment"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	commentCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, pkg.MRUsage,
	)
	listCmd.Flags().StringVarP(&parentId, "parentId", "P", "", listPidUsage)
	listCmd.Flags().StringVarP(
		&textFormat, "textFormat", "t", "html", tfUsage,
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var listTool = mcp.NewTool(
	"comment-list",
	mcp.WithTitleAnnotation(listShort),
	mcp.WithDescription(listLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(idsUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(pkg.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"parentId", mcp.DefaultString(""),
		mcp.Description(listPidUsage),
		mcp.Required(),
	),
	mcp.WithString(
		"textFormat", mcp.Enum("textFormatUnspecified", "html", "plainText"),
		mcp.DefaultString("html"), mcp.Description(tfUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet"}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(pkg.PartsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "table"),
		mcp.DefaultString("yaml"), mcp.Description(pkg.TableUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(pkg.JPUsage), mcp.Required(),
	),
)

func listHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	parentId, _ = args["parentId"].(string)
	textFormat, _ = args["textFormat"].(string)
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "comment list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "comment list failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "comment list completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithMaxResults(maxResults),
		comment.WithParentId(parentId),
		comment.WithTextFormat(textFormat),
		comment.WithService(nil),
	)

	return c.List(parts, output, jpath, writer)
}

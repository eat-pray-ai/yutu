package video

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List video's info"
	listLong     = "List video's info, such as title, description, etc"
	listIdsUsage = "Return videos with the given ids"
	listMrUsage  = "Return videos liked/disliked by the authenticated user"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	videoCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&chart, "chart", "c", "", chartUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&locale, "locale", "L", "", localUsage)
	listCmd.Flags().StringVarP(&categoryId, "videoCategoryId", "g", "", caidUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().Int64VarP(&maxHeight, "maxHeight", "H", 0, mhUsage)
	listCmd.Flags().Int64VarP(&maxWidth, "maxWidth", "W", 0, mwUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(&rating, "myRating", "R", "", listMrUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
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
	"video-list",
	mcp.WithTitleAnnotation(listShort),
	mcp.WithDescription(listLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(listIdsUsage), mcp.Required(),
	),
	mcp.WithString(
		"chart", mcp.Enum("chartUnspecified", "mostPopular"),
		mcp.DefaultString(""), mcp.Description(chartUsage), mcp.Required(),
	),
	mcp.WithString(
		"hl", mcp.DefaultString(""),
		mcp.Description(hlUsage), mcp.Required(),
	),
	mcp.WithString(
		"locale", mcp.DefaultString(""),
		mcp.Description(localUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoCategoryId", mcp.DefaultString(""),
		mcp.Description(caidUsage), mcp.Required(),
	),
	mcp.WithString(
		"regionCode", mcp.DefaultString(""),
		mcp.Description(rcUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxHeight", mcp.DefaultNumber(0),
		mcp.Description(mhUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxWidth", mcp.DefaultNumber(0),
		mcp.Description(mwUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(pkg.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"myRating", mcp.DefaultString(""),
		mcp.Description(listMrUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet", "status"}),
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
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	chart, _ = args["chart"].(string)
	hl, _ = args["hl"].(string)
	locale, _ = args["locale"].(string)
	categoryId, _ = args["videoCategoryId"].(string)
	regionCode, _ = args["regionCode"].(string)
	maxHeightRaw, _ := args["maxHeight"].(float64)
	maxHeight = int64(maxHeightRaw)
	maxWidthRaw, _ := args["maxWidth"].(float64)
	maxWidth = int64(maxWidthRaw)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	rating, _ = args["myRating"].(string)
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "video list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video list failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "video list completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithChart(chart),
		video.WithHl(hl),
		video.WithLocale(locale),
		video.WithCategory(categoryId),
		video.WithRegionCode(regionCode),
		video.WithMaxHeight(maxHeight),
		video.WithMaxWidth(maxWidth),
		video.WithMaxResults(maxResults),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithRating(rating),
		video.WithService(nil),
	)

	return v.List(parts, output, jpath, writer)
}

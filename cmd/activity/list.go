package activity

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/activity"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShortUsage = "List YouTube activities"
	listLongUsage  = "List YouTube activities, such as likes, favorites, uploads, etc."
)

var defaultParts = []string{"id", "snippet", "contentDetails"}

var listTool = mcp.NewTool(
	"activity.list",
	mcp.WithDescription(listLongUsage),
	mcp.WithString(
		"channelId", mcp.DefaultString(""), mcp.Description(channelIdUsage),
	),
	mcp.WithString(
		"home", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(homeUsage),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5), mcp.Description(maxResultsUsage),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(mineUsage),
	),
	mcp.WithString(
		"publishedAfter", mcp.DefaultString(""),
		mcp.Description(publishedAfterUsage),
	),
	mcp.WithString(
		"publishedBefore", mcp.DefaultString(""),
		mcp.Description(publishedBeforeUsage),
	),
	mcp.WithString(
		"regionCode", mcp.DefaultString(""), mcp.Description(regionCodeUsage),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray(defaultParts), mcp.Description(partsUsage),
	),
	mcp.WithString("output", mcp.DefaultString(""), mcp.Description(outputUsage)),
)

func run(writer io.Writer) error {
	a := activity.NewActivity(
		activity.WithChannelId(channelId),
		activity.WithHome(home),
		activity.WithMaxResults(maxResults),
		activity.WithMine(mine),
		activity.WithPublishedAfter(publishedAfter),
		activity.WithPublishedBefore(publishedBefore),
		activity.WithRegionCode(regionCode),
		activity.WithService(nil),
	)

	return a.List(parts, output, writer)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShortUsage,
	Long:  listLongUsage,
	Run: func(cmd *cobra.Command, args []string) {
		err := run(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", channelIdUsage,
	)
	listCmd.Flags().BoolVarP(home, "home", "H", true, homeUsage)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, maxResultsUsage,
	)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", publishedAfterUsage,
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", publishedBeforeUsage,
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", regionCodeUsage)

	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", defaultParts, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)
}

func listHandler(ctx context.Context, request mcp.CallToolRequest) (
	*mcp.CallToolResult, error,
) {
	args := request.GetArguments()
	channelId = args["channelId"].(string)
	home = utils.BoolPtr(args["home"].(string))
	maxResults = int64(args["maxResults"].(float64))
	mine = utils.BoolPtr(args["mine"].(string))
	publishedAfter = args["publishedAfter"].(string)
	publishedBefore = args["publishedBefore"].(string)
	regionCode = args["regionCode"].(string)
	parts = make([]string, len(args["parts"].([]interface{})))
	for i, part := range args["parts"].([]interface{}) {
		parts[i] = part.(string)
	}
	output = args["output"].(string)

	var writer bytes.Buffer
	err := run(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

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
	shortDesc           = "List YouTube activities"
	longDesc            = "List YouTube activities, such as likes, favorites, uploads, etc."
	channelIdDesc       = "ID of the channel"
	homeDesc            = "true or false"
	maxResultsDesc      = "The maximum number of items that should be returned"
	mineDesc            = "true or false"
	publishedAfterDesc  = "Filter on activities published after this date"
	publishedBeforeDesc = "Filter on activities published before this date"
	regionCodeDesc      = ""
	partsDesc           = "Comma separated parts"
	outputDesc          = "json or yaml"
)

var listTool = mcp.NewTool(
	"activity.list",
	mcp.WithDescription(longDesc),
	mcp.WithString(
		"channelId", mcp.DefaultString(""), mcp.Description(channelIdDesc),
	),
	mcp.WithString(
		"home", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(homeDesc),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5), mcp.Description(maxResultsDesc),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(mineDesc),
	),
	mcp.WithString(
		"publishedAfter", mcp.DefaultString(""), mcp.Description(publishedAfterDesc),
	),
	mcp.WithString(
		"publishedBefore", mcp.DefaultString(""),
		mcp.Description(publishedBeforeDesc),
	),
	mcp.WithString(
		"regionCode", mcp.DefaultString(""), mcp.Description(regionCodeDesc),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet", "contentDetails"}),
		mcp.Description(partsDesc),
	),
	mcp.WithString("output", mcp.DefaultString(""), mcp.Description(outputDesc)),
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
	Short: shortDesc,
	Long:  longDesc,
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Lookup("home").Changed {
			home = nil
		}
		if !cmd.Flags().Lookup("mine").Changed {
			mine = nil
		}
	},
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
		&channelId, "channelId", "c", "", channelIdDesc,
	)
	listCmd.Flags().BoolVarP(home, "home", "H", true, homeDesc)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, maxResultsDesc,
	)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineDesc)
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", publishedAfterDesc,
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", publishedBeforeDesc,
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", regionCodeDesc)

	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "contentDetails"}, partsDesc,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", outputDesc)
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

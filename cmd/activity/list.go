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

var defaultParts = []string{"id", "snippet", "contentDetails"}

var listTool = mcp.NewTool(
	"activity-list",
	mcp.WithTitleAnnotation("List Activities"),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithDescription(long),
	mcp.WithString(
		"channelId", mcp.DefaultString(""), mcp.Description(ciUsage), mcp.Required(),
	),
	mcp.WithString(
		"home", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(homeUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5), mcp.Description(mrUsage), mcp.Required(),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(mineUsage), mcp.Required(),
	),
	mcp.WithString(
		"publishedAfter", mcp.DefaultString(""),
		mcp.Description(paUsage), mcp.Required(),
	),
	mcp.WithString(
		"publishedBefore", mcp.DefaultString(""),
		mcp.Description(pbUsage), mcp.Required(),
	),
	mcp.WithString(
		"regionCode", mcp.DefaultString(""),
		mcp.Description(rcUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray(defaultParts),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(partsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.DefaultString("table"), mcp.Description(outputUsage),
		mcp.Required(),
	),
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
	Short: short,
	Long:  long,
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
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", ciUsage)
	listCmd.Flags().BoolVarP(home, "home", "H", true, homeUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", paUsage,
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", pbUsage,
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().StringSliceVarP(&parts, "parts", "p", defaultParts, partsUsage)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", outputUsage)
}

func listHandler(ctx context.Context, request mcp.CallToolRequest) (
	*mcp.CallToolResult, error,
) {
	args := request.GetArguments()
	channelId, _ = args["channelId"].(string)
	homeRaw, _ := args["home"].(string)
	home = utils.BoolPtr(homeRaw)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	mineRaw, ok := args["mine"].(string)
	if !ok {
		mineRaw = "true" // Default to true if not provided
	}
	mine = utils.BoolPtr(mineRaw)
	publishedAfter, _ = args["publishedAfter"].(string)
	publishedBefore, _ = args["publishedBefore"].(string)
	regionCode, _ = args["regionCode"].(string)
	partsRaw, _ := args["parts"].([]interface{})
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)

	var writer bytes.Buffer
	err := run(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

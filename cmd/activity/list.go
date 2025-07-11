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

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", ciUsage)
	listCmd.Flags().BoolVarP(home, "home", "H", true, homeUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, cmd.MRUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", paUsage,
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", pbUsage,
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", defaultParts, cmd.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var listTool = mcp.NewTool(
	"activity-list",
	mcp.WithTitleAnnotation(short),
	mcp.WithDescription(long),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(ciUsage), mcp.Required(),
	),
	mcp.WithString(
		"home", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(homeUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(mineUsage), mcp.Required(),
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
		mcp.Description(cmd.PartsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "table"),
		mcp.DefaultString("table"), mcp.Description(cmd.TableUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JPUsage), mcp.Required(),
	),
)

func listHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
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
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
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

	return a.List(parts, output, jpath, writer)
}

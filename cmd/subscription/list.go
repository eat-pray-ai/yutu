package subscription

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List subscriptions' info"
	listLong     = "List subscriptions' info, such as id, title, etc"
	listIdsUsage = "Return the subscriptions with the given ids for Stubby or Apiary"
	listCidUsage = "Return the subscriptions of the given channel owner"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	subscriptionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&forChannelId, "forChannelId", "C", "", fcidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, cmd.MRUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		myRecentSubscribers, "myRecentSubscribers", "R", false, mrsUsage,
	)
	listCmd.Flags().BoolVarP(mySubscribers, "mySubscribers", "S", false, msUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "relevance", orderUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, cmd.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)
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
	"subscription-list",
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
		"channelId", mcp.DefaultString(""),
		mcp.Description(listCidUsage), mcp.Required(),
	),
	mcp.WithString(
		"forChannelId", mcp.DefaultString(""),
		mcp.Description(fcidUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"mine", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(mineUsage), mcp.Required(),
	),
	mcp.WithString(
		"myRecentSubscribers", mcp.Enum("true", "false", ""),
		mcp.DefaultString("false"), mcp.Description(mrsUsage), mcp.Required(),
	),
	mcp.WithString(
		"mySubscribers", mcp.Enum("true", "false", ""),
		mcp.DefaultString("false"), mcp.Description(msUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwnerChannel", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"order", mcp.Enum(
			"subscriptionOrderUnspecified", "relevance", "unread", "alphabetical",
		),
		mcp.DefaultString("relevance"), mcp.Description(orderUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet"}),
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
	idsRaw, _ := args["ids"].([]any)
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	channelId, _ = args["channelId"].(string)
	forChannelId, _ = args["forChannelId"].(string)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	mineRaw, ok := args["mine"].(string)
	if !ok {
		mineRaw = "true"
	}
	mine = utils.BoolPtr(mineRaw)
	myRecentSubscribersRaw, _ := args["myRecentSubscribers"].(string)
	myRecentSubscribers = utils.BoolPtr(myRecentSubscribersRaw)
	mySubscribersRaw, _ := args["mySubscribers"].(string)
	mySubscribers = utils.BoolPtr(mySubscribersRaw)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	onBehalfOfContentOwnerChannel, _ = args["onBehalfOfContentOwnerChannel"].(string)
	order, _ = args["order"].(string)
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
	s := subscription.NewSubscription(
		subscription.WithIDs(ids),
		subscription.WithChannelId(channelId),
		subscription.WithForChannelId(forChannelId),
		subscription.WithMaxResults(maxResults),
		subscription.WithMine(mine),
		subscription.WithMyRecentSubscribers(myRecentSubscribers),
		subscription.WithMySubscribers(mySubscribers),
		subscription.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		subscription.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		subscription.WithOrder(order),
		subscription.WithService(nil),
	)

	return s.List(parts, output, jpath, writer)
}

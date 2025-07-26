package commentThread

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List YouTube comment threads"
	listLong     = "List YouTube comment threads"
	listVidUsage = "Returns the comment threads of the specified video"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	commentThreadCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().StringVarP(
		&allThreadsRelatedToChannelId, "allThreadsRelatedToChannelId", "a", "",
		atrtcidUsage,
	)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, cmd.MRUsage)
	listCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "m", "published", msUsage,
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "time", orderUsage)
	listCmd.Flags().StringVarP(&searchTerms, "searchTerms", "s", "", stUsage)
	listCmd.Flags().StringVarP(&textFormat, "textFormat", "t", "html", tfUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", listVidUsage)
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
	"commentThread-list",
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
	mcp.WithString(
		"allThreadsRelatedToChannelId", mcp.DefaultString(""),
		mcp.Description(atrtcidUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"moderationStatus",
		mcp.Enum("published", "heldForReview", "likelySpam", "rejected"),
		mcp.DefaultString("published"), mcp.Description(msUsage), mcp.Required(),
	),
	mcp.WithString(
		"order", mcp.Enum("orderUnspecified", "time", "relevance"),
		mcp.DefaultString("time"), mcp.Description(orderUsage), mcp.Required(),
	),
	mcp.WithString(
		"searchTerms", mcp.DefaultString(""),
		mcp.Description(stUsage), mcp.Required(),
	),
	mcp.WithString(
		"textFormat", mcp.Enum("textFormatUnspecified", "html"),
		mcp.DefaultString("html"), mcp.Description(tfUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoId", mcp.DefaultString(""),
		mcp.Description(listVidUsage), mcp.Required(),
	),
	mcp.WithArray(
		"parts", mcp.DefaultArray([]string{"id", "snippet"}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(cmd.PartsUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "table"),
		mcp.DefaultString("yaml"), mcp.Description(cmd.TableUsage), mcp.Required(),
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
	ids := make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	allThreadsRelatedToChannelId, _ = args["allThreadsRelatedToChannelId"].(string)
	channelId, _ = args["channelId"].(string)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	moderationStatus, _ = args["moderationStatus"].(string)
	order, _ = args["order"].(string)
	searchTerms, _ = args["searchTerms"].(string)
	textFormat, _ = args["textFormat"].(string)
	videoId, _ = args["videoId"].(string)
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
	ct := commentThread.NewCommentThread(
		commentThread.WithIDs(ids),
		commentThread.WithAllThreadsRelatedToChannelId(allThreadsRelatedToChannelId),
		commentThread.WithChannelId(channelId),
		commentThread.WithMaxResults(maxResults),
		commentThread.WithModerationStatus(moderationStatus),
		commentThread.WithOrder(order),
		commentThread.WithSearchTerms(searchTerms),
		commentThread.WithTextFormat(textFormat),
		commentThread.WithVideoId(videoId),
		commentThread.WithService(nil),
	)

	return ct.List(parts, output, jpath, writer)
}

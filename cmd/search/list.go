package search

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.MCP.AddTool(listTool, listHandler)
	searchCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&channelId, "channelId", "", cidUsage)
	listCmd.Flags().StringVar(
		&channelType, "channelType", "channelTypeUnspecified", ctUsage,
	)
	listCmd.Flags().StringVar(&eventType, "eventType", "none", etUsage)
	listCmd.Flags().BoolVar(
		forContentOwner, "forContentOwner", false, fcoUsage,
	)
	listCmd.Flags().BoolVar(forDeveloper, "forDeveloper", false, fdUsage)
	listCmd.Flags().BoolVar(forMine, "forMine", false, fmUsage)
	listCmd.Flags().StringVar(&location, "location", "", locationUsage)
	listCmd.Flags().StringVar(&locationRadius, "locationRadius", "", lrUsage)
	listCmd.Flags().Int64Var(&maxResults, "maxResults", 5, cmd.MRUsage)
	listCmd.Flags().StringVar(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "", "",
	)
	listCmd.Flags().StringVar(&order, "order", "relevance", orderUsage)
	listCmd.Flags().StringVar(&publishedAfter, "publishedAfter", "", paUsage)
	listCmd.Flags().StringVar(&publishedBefore, "publishedBefore", "", pbUsage)
	listCmd.Flags().StringVar(&q, "q", "", qUsage)
	listCmd.Flags().StringVar(&regionCode, "regionCode", "", rcUsage)
	listCmd.Flags().StringVar(&relevanceLanguage, "relevanceLanguage", "", rlUsage)
	listCmd.Flags().StringVar(&safeSearch, "safeSearch", "moderate", ssUsage)
	listCmd.Flags().StringVar(&topicId, "topicId", "", tidUsage)
	listCmd.Flags().StringSliceVar(&types, "types", []string{}, typesUsage)
	listCmd.Flags().StringVar(&videoCaption, "videoCaption", "any", vcUsage)
	listCmd.Flags().StringVar(&videoCategoryId, "videoCategoryId", "", vcidUsage)
	listCmd.Flags().StringVar(&videoDefinition, "videoDefinition", "", vdeUsage)
	listCmd.Flags().StringVar(&videoDimension, "videoDimension", "any", vdiUsage)
	listCmd.Flags().StringVar(&videoDuration, "videoDuration", "any", vduUsage)
	listCmd.Flags().StringVar(&videoEmbeddable, "videoEmbeddable", "", veUsage)
	listCmd.Flags().StringVar(&videoLicense, "videoLicense", "", vlUsage)
	listCmd.Flags().StringVar(
		&videoPaidProductPlacement, "videoPaidProductPlacement", "", vpppUsage,
	)
	listCmd.Flags().StringVar(&videoSyndicated, "videoSyndicated", "", vsUsage)
	listCmd.Flags().StringVar(&videoType, "videoType", "", vtUsage)
	listCmd.Flags().StringSliceVar(
		&parts, "parts", []string{"id", "snippet"}, cmd.PartsUsage,
	)
	listCmd.Flags().StringVar(&output, "output", "table", cmd.TableUsage)
	listCmd.Flags().StringVar(&jpath, "jsonpath", "", cmd.JPUsage)
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
	"search-list",
	mcp.WithTitleAnnotation(short),
	mcp.WithDescription(long),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelType",
		mcp.Enum("channelTypeUnspecified", "any", "show"),
		mcp.DefaultString("channelTypeUnspecified"),
		mcp.Description(ctUsage), mcp.Required(),
	),
	mcp.WithString(
		"eventType", mcp.Enum("none", "upcoming", "live", "completed"),
		mcp.DefaultString("none"), mcp.Description(etUsage), mcp.Required(),
	),
	mcp.WithString(
		"forContentOwner", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(fcoUsage), mcp.Required(),
	),
	mcp.WithString(
		"forDeveloper", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(fdUsage), mcp.Required(),
	),
	mcp.WithString(
		"forMine", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(fmUsage), mcp.Required(),
	),
	mcp.WithString(
		"location", mcp.DefaultString(""),
		mcp.Description(locationUsage), mcp.Required(),
	),
	mcp.WithString(
		"locationRadius", mcp.DefaultString(""),
		mcp.Description(lrUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"maxResults", mcp.DefaultNumber(5),
		mcp.Description(cmd.MRUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"order", mcp.DefaultString("relevance"),
		mcp.Description(orderUsage), mcp.Required(),
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
		"q", mcp.DefaultString(""),
		mcp.Description(qUsage), mcp.Required(),
	),
	mcp.WithString(
		"regionCode", mcp.DefaultString(""),
		mcp.Description(rcUsage), mcp.Required(),
	),
	mcp.WithString(
		"relevanceLanguage", mcp.DefaultString(""),
		mcp.Description(rlUsage), mcp.Required(),
	),
	mcp.WithString(
		"safeSearch",
		mcp.Enum("safeSearchSettingUnspecified", "none", "moderate", "strict"),
		mcp.DefaultString("moderate"), mcp.Description(ssUsage), mcp.Required(),
	),
	mcp.WithString(
		"topicId", mcp.DefaultString(""),
		mcp.Description(tidUsage), mcp.Required(),
	),
	mcp.WithArray(
		"types", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(typesUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoCaption",
		mcp.Enum("videoCaptionUnspecified", "any", "closedCaption", "none"),
		mcp.DefaultString("any"), mcp.Description(vcUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoCategoryId", mcp.DefaultString(""),
		mcp.Description(vcidUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoDefinition", mcp.DefaultString(""),
		mcp.Description(vdeUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoDimension", mcp.Enum("any", "2d", "3d"),
		mcp.DefaultString("any"), mcp.Description(vdiUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoDuration",
		mcp.Enum("videoDurationUnspecified", "any", "short", "medium", "long"),
		mcp.DefaultString("any"), mcp.Description(vduUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoEmbeddable", mcp.Enum("videoEmbeddableUnspecified", "any", "true"),
		mcp.DefaultString(""), mcp.Description(veUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoLicense", mcp.Enum("any", "youtube", "creativeCommon"),
		mcp.DefaultString(""), mcp.Description(vlUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoPaidProductPlacement",
		mcp.Enum("videoPaidProductPlacementUnspecified", "any", "true"),
		mcp.DefaultString(""), mcp.Description(vpppUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoSyndicated", mcp.Enum("videoSyndicatedUnspecified", "any", "true"),
		mcp.DefaultString(""), mcp.Description(vsUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoType", mcp.Enum("videoTypeUnspecified", "any", "movie", "episode"),
		mcp.DefaultString(""), mcp.Description(vtUsage), mcp.Required(),
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
	channelId, _ = args["channelId"].(string)
	channelType, _ = args["channelType"].(string)
	eventType, _ = args["eventType"].(string)
	forContentOwnerRaw, _ := args["forContentOwner"].(string)
	forContentOwner = utils.BoolPtr(forContentOwnerRaw)
	forDeveloperRaw, _ := args["forDeveloper"].(string)
	forDeveloper = utils.BoolPtr(forDeveloperRaw)
	forMineRaw, _ := args["forMine"].(string)
	forMine = utils.BoolPtr(forMineRaw)
	location, _ = args["location"].(string)
	locationRadius, _ = args["locationRadius"].(string)
	maxResultsRaw, _ := args["maxResults"].(float64)
	maxResults = int64(maxResultsRaw)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	order, _ = args["order"].(string)
	publishedAfter, _ = args["publishedAfter"].(string)
	publishedBefore, _ = args["publishedBefore"].(string)
	q, _ = args["q"].(string)
	regionCode, _ = args["regionCode"].(string)
	relevanceLanguage, _ = args["relevanceLanguage"].(string)
	safeSearch, _ = args["safeSearch"].(string)
	topicId, _ = args["topicId"].(string)
	typesRaw, _ := args["types"].([]any)
	types = make([]string, len(typesRaw))
	for i, typ := range typesRaw {
		types[i] = typ.(string)
	}
	videoCaption, _ = args["videoCaption"].(string)
	videoCategoryId, _ = args["videoCategoryId"].(string)
	videoDefinition, _ = args["videoDefinition"].(string)
	videoDimension, _ = args["videoDimension"].(string)
	videoDuration, _ = args["videoDuration"].(string)
	videoEmbeddable, _ = args["videoEmbeddable"].(string)
	videoLicense, _ = args["videoLicense"].(string)
	videoPaidProductPlacement, _ = args["videoPaidProductPlacement"].(string)
	videoSyndicated, _ = args["videoSyndicated"].(string)
	videoType, _ = args["videoType"].(string)
	partsRaw, _ := args["parts"].([]any)
	parts = make([]string, len(partsRaw))
	for i, part := range partsRaw {
		parts[i] = part.(string)
	}
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "search list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "search list failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "search list completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func list(writer io.Writer) error {
	s := search.NewSearch(
		search.WithChannelId(channelId),
		search.WithChannelType(channelType),
		search.WithEventType(eventType),
		search.WithForContentOwner(forContentOwner),
		search.WithForDeveloper(forDeveloper),
		search.WithForMine(forMine),
		search.WithLocation(location),
		search.WithLocationRadius(locationRadius),
		search.WithMaxResults(maxResults),
		search.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		search.WithOrder(order),
		search.WithPublishedAfter(publishedAfter),
		search.WithPublishedBefore(publishedBefore),
		search.WithQ(q),
		search.WithRegionCode(regionCode),
		search.WithRelevanceLanguage(relevanceLanguage),
		search.WithSafeSearch(safeSearch),
		search.WithTopicId(topicId),
		search.WithTypes(types),
		search.WithVideoCaption(videoCaption),
		search.WithVideoCategoryId(videoCategoryId),
		search.WithVideoDefinition(videoDefinition),
		search.WithVideoDimension(videoDimension),
		search.WithVideoDuration(videoDuration),
		search.WithVideoEmbeddable(videoEmbeddable),
		search.WithVideoLicense(videoLicense),
		search.WithVideoPaidProductPlacement(videoPaidProductPlacement),
		search.WithVideoSyndicated(videoSyndicated),
		search.WithVideoType(videoType),
		search.WithService(nil),
	)

	return s.List(parts, output, jpath, writer)
}

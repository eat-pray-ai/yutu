package search

import (
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List search results",
	Long:  "List search results",
	Run: func(cmd *cobra.Command, args []string) {
		s := search.NewSearch(
			search.WithChannelId(channelId),
			search.WithChannelType(channelType),
			search.WithEventType(eventType),
			search.WithForContentOwner(forContentOwner, true),
			search.WithForDeveloper(forDeveloper, true),
			search.WithForMine(forMine, true),
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
		s.List(parts, output)
	},
}

func init() {
	searchCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(
		&channelId, "channelId", "",
		"Filter on resources belonging to this channelId",
	)
	listCmd.Flags().StringVar(
		&channelType, "channelType", "",
		"channelTypeUnspecified(default), any or show",
	)
	listCmd.Flags().StringVar(
		&eventType, "eventType", "", "none(default), upcoming, live or completed",
	)
	listCmd.Flags().BoolVar(
		&forContentOwner, "forContentOwner", false, "Search owned by content owner",
	)
	listCmd.Flags().BoolVar(
		&forDeveloper, "forDeveloper", false,
		"Only retrieve videos uploaded using the project id of the authenticated user",
	)
	listCmd.Flags().BoolVar(
		&forMine, "forMine", false, "Search for the private videos of the authenticated user",
	)
	listCmd.Flags().StringVar(
		&location, "location", "", "Filter on location of the video",
	)
	listCmd.Flags().StringVar(
		&locationRadius, "locationRadius", "", "Filter on distance from the location",
	)
	listCmd.Flags().Int64Var(
		&maxResults, "maxResults", 5,
		"Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().StringVar(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "", "",
	)
	listCmd.Flags().StringVar(
		&order, "order", "",
		"searchSortUnspecified, date, rating, viewCount, relevance(default), title, videoCount",
	)
	listCmd.Flags().StringVar(
		&publishedAfter, "publishedAfter", "",
		"Filter on resources published after this date",
	)
	listCmd.Flags().StringVar(
		&publishedBefore, "publishedBefore", "",
		"Filter on resources published before this date",
	)
	listCmd.Flags().StringVar(&q, "q", "", "Textual search terms to match")
	listCmd.Flags().StringVar(
		&regionCode, "regionCode", "",
		"Display the content as seen by viewers in this country",
	)
	listCmd.Flags().StringVar(
		&relevanceLanguage, "relevanceLanguage", "",
		"Return results relevant to this language",
	)
	listCmd.Flags().StringVar(
		&safeSearch, "safeSearch", "",
		"safeSearchSettingUnspecified, none, moderate(default), strict",
	)
	listCmd.Flags().StringVar(
		&topicId, "topicId", "", "Restrict results to a particular topic",
	)
	listCmd.Flags().StringVar(
		&types, "types", "",
		"Restrict results to a particular set of resource types from One Platform",
	)
	listCmd.Flags().StringVar(
		&videoCaption, "videoCaption", "",
		"videoCaptionUnspecified, any(default), closedCaption, none",
	)
	listCmd.Flags().StringVar(
		&videoCategoryId, "videoCategoryId", "",
		"Filter on videos in a specific category",
	)
	listCmd.Flags().StringVar(
		&videoDefinition, "videoDefinition", "",
		"Filter on the definition of the videos",
	)
	listCmd.Flags().StringVar(
		&videoDimension, "videoDimension", "", "any(default), 2d or 3d",
	)
	listCmd.Flags().StringVar(
		&videoDuration, "videoDuration", "",
		"videoDurationUnspecified, any(default), short, medium, long",
	)
	listCmd.Flags().StringVar(
		&videoEmbeddable, "videoEmbeddable", "",
		"videoEmbeddableUnspecified, any or true",
	)
	listCmd.Flags().StringVar(
		&videoLicense, "videoLicense", "", "any, youtube or creativeCommon",
	)
	listCmd.Flags().StringVar(
		&videoPaidProductPlacement, "videoPaidProductPlacement", "",
		"videoPaidProductPlacementUnspecified, any or true",
	)
	listCmd.Flags().StringVar(
		&videoSyndicated, "videoSyndicated", "",
		"videoSyndicatedUnspecified, any or true",
	)
	listCmd.Flags().StringVar(
		&videoType, "videoType", "", "videoTypeUnspecified, any, movie or episode",
	)
	listCmd.Flags().StringSliceVar(
		&parts, "parts", []string{"id", "snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVar(&output, "output", "", "Output format: json or yaml")
}

package search

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List search results",
	Long:  "List search results",
	Run: func(cmd *cobra.Command, args []string) {
		s := yutuber.NewSearch(
			yutuber.WithSearchChannelId(channelId),
			yutuber.WithSearchChannelType(channelType),
			yutuber.WithSearchEventType(eventType),
			yutuber.WithSearchForContentOwner(forContentOwner),
			yutuber.WithSearchForDeveloper(forDeveloper),
			yutuber.WithSearchForMine(forMine),
			yutuber.WithSearchLocation(location),
			yutuber.WithSearchLocationRadius(locationRadius),
			yutuber.WithSearchMaxResults(maxResults),
			yutuber.WithSearchOnBehalfOfContentOwner(onBehalfOfContentOwner),
			yutuber.WithSearchOrder(order),
			yutuber.WithSearchPublishedAfter(publishedAfter),
			yutuber.WithSearchPublishedBefore(publishedBefore),
			yutuber.WithSearchQ(q),
			yutuber.WithSearchRegionCode(regionCode),
			yutuber.WithSearchRelevanceLanguage(relevanceLanguage),
			yutuber.WithSearchSafeSearch(safeSearch),
			yutuber.WithSearchTopicId(topicId),
			yutuber.WithSearchTypes(types),
			yutuber.WithSearchVideoCaption(videoCaption),
			yutuber.WithSearchVideoCategoryId(videoCategoryId),
			yutuber.WithSearchVideoDefinition(videoDefinition),
			yutuber.WithSearchVideoDimension(videoDimension),
			yutuber.WithSearchVideoDuration(videoDuration),
			yutuber.WithSearchVideoEmbeddable(videoEmbeddable),
			yutuber.WithSearchVideoLicense(videoLicense),
			yutuber.WithSearchVideoPaidProductPlacement(videoPaidProductPlacement),
			yutuber.WithSearchVideoSyndicated(videoSyndicated),
			yutuber.WithSearchVideoType(videoType),
		)
		s.List(parts, output)
	},
}

func init() {
	searchCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&channelId, "channelId", "", "channel id")
	listCmd.Flags().StringVar(
		&channelType, "channelType", "", "channelTypeUnspecified(default), any or show",
	)
	listCmd.Flags().StringVar(&eventType, "eventType", "", "none(default), upcoming, live or completed")
	listCmd.Flags().BoolVar(&forContentOwner, "forContentOwner", false, "search owned by content owner")
	listCmd.Flags().BoolVar(
		&forDeveloper, "forDeveloper", false,
		"only retrieve videos uploaded using the project id of the authenticated user",
	)
	listCmd.Flags().BoolVar(&forMine, "forMine", false, "search for the private videos of the authenticated user")
	listCmd.Flags().StringVar(&location, "location", "", "filter on location of the video")
	listCmd.Flags().StringVar(&locationRadius, "locationRadius", "", "filter on distance from the location")
	listCmd.Flags().Int64Var(
		&maxResults, "maxResults", 5,
		"specifies the maximum number of items that should be returned in the result set",
	)
	listCmd.Flags().StringVar(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "", "")
	listCmd.Flags().StringVar(
		&order, "order", "",
		"searchSortUnspecified, date, rating, viewCount, relevance(default), title, videoCount",
	)
	listCmd.Flags().StringVar(&publishedAfter, "publishedAfter", "", "filter on resources published after this date")
	listCmd.Flags().StringVar(&publishedBefore, "publishedBefore", "", "filter on resources published before this date")
	listCmd.Flags().StringVar(&q, "q", "", "textual search terms to match")
	listCmd.Flags().StringVar(&regionCode, "regionCode", "", "display the content as seen by viewers in this country")
	listCmd.Flags().StringVar(&relevanceLanguage, "relevanceLanguage", "", "return results relevant to this language")
	listCmd.Flags().StringVar(
		&safeSearch, "safeSearch", "",
		"safeSearchSettingUnspecified, none, moderate(default), strict",
	)
	listCmd.Flags().StringVar(&topicId, "topicId", "", "restrict results to a particular topic")
	listCmd.Flags().StringVar(
		&types, "types", "", "restrict results to a particular set of resource types from One Platform",
	)
	listCmd.Flags().StringVar(
		&videoCaption, "videoCaption", "",
		"videoCaptionUnspecified, any(default), closedCaption, none",
	)
	listCmd.Flags().StringVar(&videoCategoryId, "videoCategoryId", "", "filter on videos in a specific category")
	listCmd.Flags().StringVar(&videoDefinition, "videoDefinition", "", "filter on the definition of the videos")
	listCmd.Flags().StringVar(&videoDimension, "videoDimension", "", "any(default), 2d or 3d")
	listCmd.Flags().StringVar(
		&videoDuration, "videoDuration", "",
		"videoDurationUnspecified, any(default), short, medium, long",
	)
	listCmd.Flags().StringVar(&videoEmbeddable, "videoEmbeddable", "", "videoEmbeddableUnspecified, any or true")
	listCmd.Flags().StringVar(&videoLicense, "videoLicense", "", "any, youtube or creativeCommon")
	listCmd.Flags().StringVar(
		&videoPaidProductPlacement, "videoPaidProductPlacement", "",
		"videoPaidProductPlacementUnspecified, any or true",
	)
	listCmd.Flags().StringVar(&videoSyndicated, "videoSyndicated", "", "videoSyndicatedUnspecified, any or true")
	listCmd.Flags().StringVar(&videoType, "videoType", "", "videoTypeUnspecified, any, movie or episode")
	listCmd.Flags().StringSliceVar(&parts, "parts", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVar(&output, "output", "", "Output format: json or yaml")
}

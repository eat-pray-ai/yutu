package search

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/spf13/cobra"
	"io"
)

func init() {
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
	listCmd.Flags().Int64Var(&maxResults, "maxResults", 5, mrUsage)
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
		&parts, "parts", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVar(&output, "output", "table", cmd.TableUsage)
	listCmd.Flags().StringVar(&jpath, "jsonpath", "", cmd.JpUsage)
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

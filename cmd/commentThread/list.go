package commentThread

import (
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List YouTube comment threads",
	Long:  "List YouTube comment threads",
	Run: func(cmd *cobra.Command, args []string) {
		ct := commentThread.NewCommentThread(
			commentThread.WithID(id),
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
		ct.List(parts, output)
	},
}

func init() {
	commentThreadCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&id, "id", "i", []string{}, "ID of the comment thread")
	listCmd.Flags().StringVarP(
		&allThreadsRelatedToChannelId, "allThreadsRelatedToChannelId", "a", "",
		"Returns the comment threads of all videos of the channel and the channel comments as well",
	)
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", "Returns the comment threads for all the channel comments",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, "Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "m", "", "published(default), heldForReview, likelySpam or rejected",
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "", "orderUnspecified, time(default) or relevance")
	listCmd.Flags().StringVarP(&searchTerms, "searchTerms", "s", "", "Search terms")
	listCmd.Flags().StringVarP(&textFormat, "textFormat", "t", "", "textFormatUnspecified or html(default)")
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "Returns the comment threads of the specified video")
	listCmd.Flags().StringSliceVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Parts to be fetched")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

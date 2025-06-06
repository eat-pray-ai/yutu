package commentThread

import (
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List YouTube comment threads"
	listLong     = "List YouTube comment threads"
	listVidUsage = "Returns the comment threads of the specified video"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
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
		ct.List(parts, output)
	},
}

func init() {
	commentThreadCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().StringVarP(
		&allThreadsRelatedToChannelId, "allThreadsRelatedToChannelId", "a", "",
		atrtcidUsage,
	)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "m", "published", msUsage,
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "time", orderUsage)
	listCmd.Flags().StringVarP(&searchTerms, "searchTerms", "s", "", stUsage)
	listCmd.Flags().StringVarP(&textFormat, "textFormat", "t", "html", tfUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", listVidUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)
}

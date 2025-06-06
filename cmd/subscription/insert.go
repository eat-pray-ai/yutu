package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
)

const (
	insertShort       = "Insert a YouTube subscription"
	insertLong        = "Insert a YouTube subscription"
	insertCidUsage    = "ID of the channel to be subscribed"
	insertOutputUsage = "json, yaml, or silent"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithSubscriberChannelId(subscriberChannelId),
			subscription.WithDescription(description),
			subscription.WithChannelId(channelId),
			subscription.WithTitle(title),
			subscription.WithService(nil),
		)
		s.Insert(output)
	},
}

func init() {
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", insertOutputUsage)

	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

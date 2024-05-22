package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a subscription",
	Long:  "Insert a subscription",
	Run: func(cmd *cobra.Command, args []string) {
		s := yutuber.NewSubscription(
			yutuber.WithSubscriptionSubscriberChannelId(subscriberChannelId),
			yutuber.WithSubscriptionDescription(description),
			yutuber.WithSubscriptionChannelId(channelId),
			yutuber.WithSubscriptionTitle(title),
		)
		s.Insert()
	},
}

func init() {
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&subscriberChannelId, "subscriberChannelId", "s", "", "Subscriber's channel ID")
	insertCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the subscription")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of the channel to be subscribed")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the subscription")
}

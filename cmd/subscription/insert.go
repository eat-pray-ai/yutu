package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/subscription"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a subscription",
	Long:  "Insert a subscription",
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
		&subscriberChannelId, "subscriberChannelId", "s", "",
		"Subscriber's channel ID",
	)
	insertCmd.Flags().StringVarP(
		&description, "description", "d", "", "Description of the subscription",
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", "ID of the channel to be subscribed",
	)
	insertCmd.Flags().StringVarP(
		&title, "title", "t", "", "Title of the subscription",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

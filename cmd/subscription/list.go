package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list subscriptions' info",
	Long:  "list subscriptions' info, such as id, title, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		s := yutuber.NewSubscription(
			yutuber.WithSubscriptionId(id),
			yutuber.WithSubscriptionChannelId(channelId),
			yutuber.WithSubscriptionForChannelId(forChannelId),
			yutuber.WithSubscriptionMaxResults(maxResults),
			yutuber.WithSubscriptionMine(mine),
			yutuber.WithSubscriptionMyRecentSubscribers(myRecentSubscribers),
			yutuber.WithSubscriptionMySubscribers(mySubscribers),
			yutuber.WithSubscriptionOnBehalfOfContentOwner(onBehalfOfContentOwner),
			yutuber.WithSubscriptionOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			yutuber.WithSubscriptionOrder(order),
		)
		s.List(parts, output)
	},
}

func init() {
	subscriptionCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&id, "id", "i", "",
		"Return the subscriptions with the given IDs for Stubby or Apiary",
	)
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "",
		"Return the subscriptions of the given channel owner",
	)
	listCmd.Flags().StringVarP(
		&forChannelId, "forChannelId", "C", "",
		"Return the subscriptions to the subset of these channels that the authenticated user is subscribed to",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5,
		"Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(&mine, "mine", "m", "", "true or false")
	listCmd.Flags().StringVarP(&myRecentSubscribers, "myRecentSubscribers", "r", "", "true  or false")
	listCmd.Flags().StringVarP(&mySubscribers, "mySubscribers", "s", "", "true  or false")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "")
	listCmd.Flags().StringVarP(
		&order, "order", "O", "", "subscriptionOrderUnspecified, relevance(default), unread or alphabetical",
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
}

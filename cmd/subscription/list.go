package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List subscriptions' info",
	Long:  "List subscriptions' info, such as id, title, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
			subscription.WithID(id),
			subscription.WithChannelId(channelId),
			subscription.WithForChannelId(forChannelId),
			subscription.WithMaxResults(maxResults),
			subscription.WithMine(mine, true),
			subscription.WithMyRecentSubscribers(myRecentSubscribers, true),
			subscription.WithMySubscribers(mySubscribers, true),
			subscription.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			subscription.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			subscription.WithOrder(order),
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
		&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned",
	)
	listCmd.Flags().BoolVarP(&mine, "mine", "M", true, "Return the subscriptions of the authenticated user")
	listCmd.Flags().BoolVarP(
		&myRecentSubscribers, "myRecentSubscribers", "R", false, "true  or false",
	)
	listCmd.Flags().BoolVarP(
		&mySubscribers, "mySubscribers", "S", false, "Return the subscribers of the given channel owner",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringVarP(
		&order, "order", "O", "",
		"subscriptionOrderUnspecified, relevance(default), unread or alphabetical",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "json or yaml",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts",
	)
}

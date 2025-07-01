package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List subscriptions' info"
	listLong     = "List subscriptions' info, such as id, title, etc"
	listIdsUsage = "Return the subscriptions with the given ids for Stubby or Apiary"
	listCidUsage = "Return the subscriptions of the given channel owner"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithIDs(ids),
			subscription.WithChannelId(channelId),
			subscription.WithForChannelId(forChannelId),
			subscription.WithMaxResults(maxResults),
			subscription.WithMine(mine),
			subscription.WithMyRecentSubscribers(myRecentSubscribers),
			subscription.WithMySubscribers(mySubscribers),
			subscription.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			subscription.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			subscription.WithOrder(order),
			subscription.WithService(nil),
		)

		err := s.List(parts, output, jpath, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	subscriptionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&forChannelId, "forChannelId", "C", "", fcidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		myRecentSubscribers, "myRecentSubscribers", "R", false, mrsUsage,
	)
	listCmd.Flags().BoolVarP(mySubscribers, "mySubscribers", "S", false, msUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "relevance", orderUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)
}

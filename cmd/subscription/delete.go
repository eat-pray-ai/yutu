package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/subscription"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete subscription",
	Long:  "Delete subscription",
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithID(id), subscription.WithService(nil),
		)
		s.Delete()
	},
}

func init() {
	subscriptionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(
		&id, "id", "i", "", "ID of the subscription to delete",
	)
}

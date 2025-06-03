package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
)

const (
	deleteShort   = "Delete a YouTube subscription"
	deleteLong    = "Delete a YouTube subscription by id"
	deleteIdUsage = "ID of the subscription to delete"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithID(id), subscription.WithService(nil),
		)
		s.Delete()
	},
}

func init() {
	subscriptionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", deleteIdUsage)
}

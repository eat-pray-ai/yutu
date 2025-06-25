package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete a YouTube subscriptions"
	deleteLong     = "Delete a YouTube subscriptions by ids"
	deleteIdsUsage = "IDs of the subscriptions to delete"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		s := subscription.NewSubscription(
			subscription.WithIDs(ids), subscription.WithService(nil),
		)

		err := s.Delete(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	subscriptionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

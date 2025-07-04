package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort    = "Delete a YouTube subscriptions"
	deleteLong     = "Delete a YouTube subscriptions by ids"
	deleteIdsUsage = "IDs of the subscriptions to delete"
)

func init() {
	subscriptionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func del(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithIDs(ids), subscription.WithService(nil),
	)

	return s.Delete(writer)
}

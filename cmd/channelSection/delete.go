package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete channel sections"
	deleteLong     = "Delete channel sections by ids"
	deleteIdsUsage = "Delete the channel sections with the given ids"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		cs := channelSection.NewChannelSection(
			channelSection.WithIDs(ids),
			channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channelSection.WithService(nil),
		)
		cs.Delete()
	},
}

func init() {
	channelSectionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
}

package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/spf13/cobra"
)

const (
	deleteShort   = "Delete channel section"
	deleteLong    = "Delete channel section"
	deleteIdUsage = "Delete the channel section with the given id"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		cs := channelSection.NewChannelSection(
			channelSection.WithID(id),
			channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channelSection.WithService(nil),
		)
		cs.Delete()
	},
}

func init() {
	channelSectionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", deleteIdUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
}

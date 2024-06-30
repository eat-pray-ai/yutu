package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/channelSection"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete channel section",
	Long:  "Delete channel section",
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

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", "Delete the ChannelSections with the given ID")
	deleteCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
}

package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/spf13/cobra"
)

const (
	listShort   = "List channel sections"
	listLong    = "List channel sections"
	listIdUsage = "Return the channel sections with the given id"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		cs := channelSection.NewChannelSection(
			channelSection.WithID(id), // todo: id -> ids
			channelSection.WithChannelId(channelId),
			channelSection.WithHl(hl),
			channelSection.WithMine(mine),
			channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channelSection.WithService(nil),
		)
		cs.List(parts, output)
	},
}

func init() {
	channelSectionCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", listIdUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", false, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "json", outputUsage)
}

package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/channelSection"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List channel sections",
	Long:  "List channel sections",
	Run: func(cmd *cobra.Command, args []string) {
		cs := channelSection.NewChannelSection(
			channelSection.WithId(id),
			channelSection.WithChannelId(channelId),
			channelSection.WithHl(hl),
			channelSection.WithMine(mine, true),
			channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		cs.List(parts, output)
	},
}

func init() {
	channelSectionCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "Return the ChannelSections with the given ID")
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", " Return the ChannelSections owned by the specified channel ID",
	)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", "Return content in specified language")
	listCmd.Flags().BoolVarP(&mine, "mine", "m", false, "Return the ChannelSections owned by the authenticated user")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "json", "json or yaml")
}

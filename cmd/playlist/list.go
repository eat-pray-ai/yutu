package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
)

const (
	listShort       = "List playlist's info"
	listLong        = "List playlist's info, such as title, description, etc"
	listIdUsage     = "Return the playlists with the given IDs for Stubby or Apiary"
	listCidUsage    = "Return the playlists owned by the specified channel id"
	listOutputUsage = "json or yaml"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithID(id), // todo: id -> ids
			playlist.WithChannelId(channelId),
			playlist.WithHl(hl),
			playlist.WithMaxResults(maxResults),
			playlist.WithMine(mine),
			playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlist.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			playlist.WithService(nil),
		)
		p.List(parts, output)
	},
}

func init() {
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", listIdUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
}

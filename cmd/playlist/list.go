package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlist"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List playlist's info",
	Long:  "List playlist's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithID(id),
			playlist.WithChannelId(channelId),
			playlist.WithHl(hl),
			playlist.WithMaxResults(maxResults),
			playlist.WithMine(mine, true),
			playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlist.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			playlist.WithService(nil),
		)
		p.List(parts, output)
	},
}

func init() {
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&id, "id", "i", "",
		"Return the playlists with the given IDs for Stubby or Apiary.",
	)
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "",
		"Return the playlists owned by the specified channel ID",
	)
	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Return content in specified language",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5,
		"Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().BoolVarP(&mine, "mine", "m", true, "Return the playlists owned by the authenticated user")
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"},
		"Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

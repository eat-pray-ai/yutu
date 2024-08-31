package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List playlist items",
	Long:  "List playlist items' info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithID(id),
			playlistItem.WithPlaylistId(playlistId),
			playlistItem.WithMaxResults(maxResults),
			playlistItem.WithVideoId(videoId),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithService(nil),
		)
		pi.List(parts, output)
	},
}

func init() {
	playlistItemCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the playlist item")
	listCmd.Flags().StringVarP(
		&playlistId, "playlistId", "y", "",
		"Return the playlist items within the given playlist",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(
		&videoId, "videoId", "v", "",
		"Return the playlist items associated with the given video ID",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "json or yaml",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"},
		"Comma separated parts",
	)
}

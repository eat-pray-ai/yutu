package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
)

const (
	listShort       = "List playlist items"
	listLong        = "List playlist items' info, such as title, description, etc"
	listIdUsage     = "ID of the playlist item to list"
	listPidUsage    = "Return the playlist items within the given playlist"
	listOutputUsage = "json or yaml"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithID(id), // todo: id -> ids
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

	listCmd.Flags().StringVarP(&id, "id", "i", "", listIdUsage)
	listCmd.Flags().StringVarP(&playlistId, "playlistId", "y", "", listPidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
}

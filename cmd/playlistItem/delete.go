package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlistItem"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a item from a playlist",
	Long:  "delete a item from a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithId(id),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithService(),
		)
		pi.Delete()
	},
}

func init() {
	playlistItemCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(
		&id, "id", "i", "", "ID of the playlist item to be deleted",
	)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
}

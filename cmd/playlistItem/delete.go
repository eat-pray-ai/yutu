package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a item from a playlist",
	Long:  "Delete a item from a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithID(id),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithService(nil),
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

package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a playlist",
	Long:  "delete a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		p := yutuber.NewPlaylist(
			yutuber.WithPlaylistId(id),
			yutuber.WithPlaylistOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		p.Delete()
	},
}

func init() {
	playlistCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", "ID of playlist to be deleted")
	deleteCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
}

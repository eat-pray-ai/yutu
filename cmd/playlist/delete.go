package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a playlist",
	Long:  "Delete a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithID(id),
			playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlist.WithService(nil),
		)
		p.Delete()
	},
}

func init() {
	playlistCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(
		&id, "id", "i", "", "ID of playlist to be deleted",
	)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
}

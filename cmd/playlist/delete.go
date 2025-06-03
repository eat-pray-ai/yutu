package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
)

const (
	deleteShort   = "Delete a playlist"
	deleteLong    = "Delete a playlist by id"
	deleteIdUsage = "ID of the playlist to delete"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
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

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", deleteIdUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
}

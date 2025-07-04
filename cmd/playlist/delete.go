package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort    = "Delete a playlists"
	deleteLong     = "Delete a playlists by ids"
	deleteIdsUsage = "IDs of the playlists to delete"
)

func init() {
	playlistCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func del(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithIDs(ids),
		playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlist.WithService(nil),
	)

	return p.Delete(writer)
}

package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort    = "Delete items from a playlist"
	deleteLong     = "Delete items from a playlist by ids"
	deleteIdsUsage = "IDs of the playlist items to delete"
)

func init() {
	playlistItemCmd.AddCommand(deleteCmd)

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
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithIDs(ids),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.Delete(writer)
}

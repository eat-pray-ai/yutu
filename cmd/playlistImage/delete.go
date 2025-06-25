package playlistImage

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/spf13/cobra"
)

const (
	deleteShort = "Delete YouTube playlist images"
	deleteLong  = "Delete YouTube playlist images by ids"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithIDs(ids),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithService(nil),
		)

		err := pi.Delete(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	playlistImageCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = deleteCmd.MarkFlagRequired("ids")
}

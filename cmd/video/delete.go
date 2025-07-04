package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort    = "Delete a video on YouTube"
	deleteLong     = "Delete a video on YouTube by ids"
	deleteIdsUsage = "IDs of the videos to delete"
)

func init() {
	videoCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
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
	v := video.NewVideo(video.WithIDs(ids))

	return v.Delete(writer)
}

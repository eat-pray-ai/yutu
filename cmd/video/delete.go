package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete a video on YouTube"
	deleteLong     = "Delete a video on YouTube by ids"
	deleteIdsUsage = "IDs of the videos to delete"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(video.WithIDs(ids))
		v.Delete()
	},
}

func init() {
	videoCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

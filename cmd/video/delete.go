package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	deleteShort   = "Delete a video on YouTube"
	deleteLong    = "Delete a video on YouTube by id"
	deleteIdUsage = "ID of the video to delete"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(video.WithID(id))
		v.Delete()
	},
}

func init() {
	videoCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", deleteIdUsage)
	deleteCmd.MarkFlagRequired("id")
}

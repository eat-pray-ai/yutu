package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/video"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a video on Youtube",
	Long:  "Delete a video on Youtube",
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(video.WithId(id))
		v.Delete()
	},
}

func init() {
	videoCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video to delete")
	deleteCmd.MarkFlagRequired("id")
}

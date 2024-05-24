package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var getRatingCmd = &cobra.Command{
	Use:   "getRating",
	Short: "get the rating of a video",
	Long:  "get the rating of a video, with the specified video ID",
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithVideoId(id),
			yutuber.WithVideoOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		v.GetRating()
	},
}

func init() {
	videoCmd.AddCommand(getRatingCmd)

	getRatingCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video")
	getRatingCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	getRatingCmd.MarkFlagRequired("id")
}

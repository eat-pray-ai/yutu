package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: "Rate a video on YouTube",
	Long:  "Rate a video on YouTube, with the specified rating",
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id),
			video.WithRating(rating),
			video.WithService(nil),
		)
		v.Rate()
	},
}

func init() {
	videoCmd.AddCommand(rateCmd)

	rateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video")
	rateCmd.Flags().StringVarP(
		&rating, "rating", "r", "", "Rating of the video: like, dislike or none",
	)

	rateCmd.MarkFlagRequired("id")
	rateCmd.MarkFlagRequired("rating")
}

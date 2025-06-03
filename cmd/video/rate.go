package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	rateShort   = "Rate a video on YouTube"
	rateLong    = "Rate a video on YouTube, with the specified rating"
	rateIdUsage = "ID of the video to rate"
	rateRUsage  = "like, dislike, or none"
)

var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: rateShort,
	Long:  rateLong,
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

	rateCmd.Flags().StringVarP(&id, "id", "i", "", rateIdUsage)
	rateCmd.Flags().StringVarP(&rating, "rating", "r", "", rateRUsage)

	rateCmd.MarkFlagRequired("id")
	rateCmd.MarkFlagRequired("rating")
}

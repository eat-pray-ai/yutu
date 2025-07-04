package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
	"io"
)

const (
	rateShort    = "Rate a video on YouTube"
	rateLong     = "Rate a video on YouTube, with the specified rating"
	rateIdsUsage = "IDs of the videos to rate"
	rateRUsage   = "like, dislike, or none"
)

func init() {
	videoCmd.AddCommand(rateCmd)

	rateCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, rateIdsUsage)
	rateCmd.Flags().StringVarP(&rating, "rating", "r", "", rateRUsage)

	_ = rateCmd.MarkFlagRequired("ids")
	_ = rateCmd.MarkFlagRequired("rating")
}

var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: rateShort,
	Long:  rateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := rate(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func rate(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithRating(rating),
		video.WithService(nil),
	)

	return v.Rate(writer)
}

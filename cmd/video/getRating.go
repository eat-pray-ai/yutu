package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	getRatingShort = "Get the rating of a video"
	getRatingLong  = "Get the rating of a video by id"
	grIdUsage      = "ID of the video"
)

var getRatingCmd = &cobra.Command{
	Use:   "getRating",
	Short: getRatingShort,
	Long:  getRatingLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithService(nil),
		)
		v.GetRating()
	},
}

func init() {
	videoCmd.AddCommand(getRatingCmd)

	getRatingCmd.Flags().StringVarP(&id, "id", "i", "", grIdUsage)
	getRatingCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	getRatingCmd.MarkFlagRequired("id")
}

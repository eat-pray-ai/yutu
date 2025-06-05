package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	getRatingShort = "Get the rating of videos"
	getRatingLong  = "Get the rating of videos by ids"
	grIdsUsage     = "IDs of the videos to get the rating for"
	grOutputUsage  = "json or yaml"
)

var getRatingCmd = &cobra.Command{
	Use:   "getRating",
	Short: getRatingShort,
	Long:  getRatingLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithIDs(ids),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithService(nil),
		)
		v.GetRating(output)
	},
}

func init() {
	videoCmd.AddCommand(getRatingCmd)

	getRatingCmd.Flags().StringArrayVarP(&ids, "ids", "i", []string{}, grIdsUsage)
	getRatingCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	getRatingCmd.Flags().StringVarP(&output, "output", "o", "", grOutputUsage)
	getRatingCmd.MarkFlagRequired("id")
}

package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List video's info",
	Long:  "List video's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id),
			video.WithChart(chart),
			video.WithHl(hl),
			video.WithLocale(locale),
			video.WithCategory(categoryId),
			video.WithRegionCode(regionCode),
			video.WithMaxHeight(maxHeight),
			video.WithMaxWidth(maxWidth),
			video.WithMaxResults(maxResults),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithRating(rating),
			video.WithService(nil),
		)
		v.List(parts, output)
	},
}

func init() {
	videoCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&id, "id", "i", "", "Return videos with the given ids",
	)
	listCmd.Flags().StringVarP(
		&chart, "chart", "c", "", "chartUnspecified or mostPopular",
	)
	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Specifies the localization language",
	)
	listCmd.Flags().StringVarP(&locale, "locale", "L", "", "")
	listCmd.Flags().StringVarP(
		&categoryId, "videoCategoryId", "g", "",
		"Specific to the specified video category",
	)
	listCmd.Flags().StringVarP(
		&regionCode, "regionCode", "r", "", "Specific to the specified region",
	)
	listCmd.Flags().Int64VarP(&maxHeight, "maxHeight", "H", 0, "")
	listCmd.Flags().Int64VarP(&maxWidth, "maxWidth", "W", 0, "")
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5,
		"Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&rating, "myRating", "R", "",
		"Return videos liked/disliked by the authenticated user",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "json or yaml",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, "Comma separated parts",
	)
}

package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	listShort       = "List video's info"
	listLong        = "List video's info, such as title, description, etc"
	listIdUsage     = "Return videos with the given ids"
	listMrUsage     = "Return videos liked/disliked by the authenticated user"
	listOutputUsage = "json or yaml"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id), // todo: id -> ids
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

	listCmd.Flags().StringVarP(&id, "id", "i", "", listIdUsage)
	listCmd.Flags().StringVarP(&chart, "chart", "c", "", chartUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&locale, "locale", "L", "", localUsage)
	listCmd.Flags().StringVarP(&categoryId, "videoCategoryId", "g", "", caidUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().Int64VarP(&maxHeight, "maxHeight", "H", 0, mhUsage)
	listCmd.Flags().Int64VarP(&maxWidth, "maxWidth", "W", 0, mwUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(&rating, "myRating", "R", "", listMrUsage)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
}

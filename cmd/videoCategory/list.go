package videoCategory

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/videoCategory"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list video categories",
	Long:  "list video categories' info, such as ID, title, assignable, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		vc := videoCategory.NewVideoCategory(
			videoCategory.WithId(id),
			videoCategory.WithHl(hl),
			videoCategory.WithRegionCode(regionCode),
		)
		vc.List(parts, output)
	},
}

func init() {
	videoCategoryCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video category")
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", "Host language")
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "US", "Region code")
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

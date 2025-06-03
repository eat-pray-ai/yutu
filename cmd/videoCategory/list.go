package videoCategory

import (
	"github.com/eat-pray-ai/yutu/pkg/videoCategory"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		vc := videoCategory.NewVideoCategory(
			videoCategory.WithID(id),
			videoCategory.WithHl(hl),
			videoCategory.WithRegionCode(regionCode),
			videoCategory.WithService(nil),
		)
		vc.List(parts, output)
	},
}

func init() {
	videoCategoryCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", idUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "US", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)
}

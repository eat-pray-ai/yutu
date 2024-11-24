package videoCategory

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/videoCategory"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List video categories",
	Long:  "List video categories' info, such as ID, title, assignable, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		vc := videoCategory.NewVideoCategory(
			videoCategory.WithID(id),
			videoCategory.WithHl(hl),
			videoCategory.WithRegionCode(regionCode),
			videoCategory.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
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
		&output, "output", "o", "", "json or yaml",
	)
}

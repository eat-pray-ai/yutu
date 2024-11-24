package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/i18nRegion"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List i18nRegions",
	Long:  "List i18nRegions' id, hl, and name",
	Run: func(cmd *cobra.Command, args []string) {
		i := i18nRegion.NewI18nRegion(
			i18nRegion.WithHl(hl),
			i18nRegion.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		i.List(parts, output)
	},
}

func init() {
	i18nRegionCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", "Host language")
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "json or yaml",
	)
}

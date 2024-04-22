package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list i18nRegions",
	Long:  "list i18nRegions' id, hl, and name",
	Run: func(cmd *cobra.Command, args []string) {
		i := yutuber.NewI18nRegion()
		i.List(parts, output)
	},
}

func init() {
	i18nRegionCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "output format: json or yaml",
	)
}

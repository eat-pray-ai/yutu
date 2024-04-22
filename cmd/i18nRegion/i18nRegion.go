package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	parts  []string
	output string
)

var i18nRegionCmd = &cobra.Command{
	Use:   "i18nRegion",
	Short: "manipulate YouTube i18nRegions",
	Long:  "manipulate YouTube i18nRegions, only list for now",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nRegionCmd)
}

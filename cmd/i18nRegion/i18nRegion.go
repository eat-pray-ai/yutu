package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	hl     string
	parts  []string
	output string
)

var i18nRegionCmd = &cobra.Command{
	Use:   "i18nRegion",
	Short: "List YouTube i18n regions",
	Long:  "List YouTube i18n regions",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nRegionCmd)
}

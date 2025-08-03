package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short   = "List YouTube i18n regions"
	long    = "List YouTube i18n regions' id, hl, and name"
	hlUsage = "Host language"
	rURI    = "i18n://region"
	rName   = "all i18n region"
)

var (
	hl           string
	parts        []string
	output       string
	jpath        string
	defaultParts = []string{"id", "snippet"}
)

var i18nRegionCmd = &cobra.Command{
	Use:   "i18nRegion",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nRegionCmd)
}

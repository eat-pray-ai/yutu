package i18nLanguage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short    = "List YouTube i18n languages"
	long     = "List YouTube i18n languages' id, hl, and name"
	hlUsage  = "Host language"
	hlURI    = "i18n://hl"
	hlName   = "all i18n host languages"
	hlDesc   = "List all i18n host languages for YouTube regions"
	langURI  = "i18n://language"
	langName = "all i18n languages"
)

var (
	hl           string
	parts        []string
	output       string
	jpath        string
	defaultParts = []string{"id", "snippet"}
)

var i18nLanguageCmd = &cobra.Command{
	Use:   "i18nLanguage",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nLanguageCmd)
}

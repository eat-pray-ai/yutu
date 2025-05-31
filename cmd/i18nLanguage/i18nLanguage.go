package i18nLanguage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	hl     string
	parts  []string
	output string
)

var i18nLanguageCmd = &cobra.Command{
	Use:   "i18nLanguage",
	Short: "List YouTube i18n languages",
	Long:  "List YouTube i18n languages",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nLanguageCmd)
}

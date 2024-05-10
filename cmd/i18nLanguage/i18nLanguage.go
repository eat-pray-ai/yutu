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
	Short: "manipulate YouTube i18nLanguages",
	Long:  "manipulate YouTube i18nLanguages, only list for now",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nLanguageCmd)
}

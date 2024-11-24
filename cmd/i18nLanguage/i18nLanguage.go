package i18nLanguage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	hl         string
	parts      []string
	output     string
	credential string
	cacheToken string
)

var i18nLanguageCmd = &cobra.Command{
	Use:   "i18nLanguage",
	Short: "List YouTube i18nLanguages",
	Long:  "List YouTube i18nLanguages",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nLanguageCmd)

	i18nLanguageCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	i18nLanguageCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

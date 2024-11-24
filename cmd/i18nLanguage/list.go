package i18nLanguage

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/i18nLanguage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List i18nLanguages",
	Long:  "List i18nLanguages' id, hl, and name",
	Run: func(cmd *cobra.Command, args []string) {
		i := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithHl(hl),
			i18nLanguage.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		i.List(parts, output)
	},
}

func init() {
	i18nLanguageCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", "Host language")
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "json or yaml",
	)
}

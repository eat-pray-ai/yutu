package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id         string
	hl         string
	regionCode string
	parts      []string
	output     string
	credential string
	cacheToken string
)

var videoCategoryCmd = &cobra.Command{
	Use:   "videoCategory",
	Short: "List YouTube video categories",
	Long:  "List YouTube video categories",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCategoryCmd)

	videoCategoryCmd.PersistentFlags().StringVarP(
		&credential, "credential", "c", "client_secret.json", "Path to client secret file",
	)
	videoCategoryCmd.PersistentFlags().StringVarP(
		&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file",
	)
}

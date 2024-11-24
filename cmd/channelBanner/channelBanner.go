package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	file                          string
	output                        string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	credential                    string
	cacheToken                    string
)

var channelBannerCmd = &cobra.Command{
	Use:   "channelBanner",
	Short: "Insert Youtube channelBanner",
	Long:  "Insert Youtube channelBanner",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelBannerCmd)

	channelBannerCmd.PersistentFlags().StringVarP(
		&credential, "credential", "c", "client_secret.json", "Path to client secret file",
	)
	channelBannerCmd.PersistentFlags().StringVarP(
		&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file",
	)
}

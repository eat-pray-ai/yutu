package member

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	memberChannelId  string
	hasAccessToLevel string
	maxResults       int64
	mode             string
	parts            []string
	output           string
	credential       string
	cacheToken       string
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "List YouTube members",
	Long:  "List YouTube members",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(memberCmd)

	memberCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	memberCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

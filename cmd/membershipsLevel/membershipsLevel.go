package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	parts      []string
	output     string
	credential string
	cacheToken string
)

var membershipsLevelCmd = &cobra.Command{
	Use:   "membershipsLevel",
	Short: "List YouTube memberships levels",
	Long:  "List YouTube memberships levels",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)

	membershipsLevelCmd.PersistentFlags().StringVarP(
		&credential, "credential", "c", "client_secret.json", "Path to client secret file",
	)
	membershipsLevelCmd.PersistentFlags().StringVarP(
		&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file",
	)
}

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
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "List YouTube members",
	Long:  "List YouTube members",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(memberCmd)
}

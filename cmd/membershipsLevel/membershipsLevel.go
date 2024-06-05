package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	parts  []string
	output string
)

var membershipsLevelCmd = &cobra.Command{
	Use:   "membershipsLevel",
	Short: "list YouTube memberships levels",
	Long:  "list YouTube memberships levels",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)
}

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
	Short: "List YouTube memberships levels",
	Long:  "List YouTube memberships levels",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)
}

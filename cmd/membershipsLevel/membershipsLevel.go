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
	Short: "manipulate YouTube memberships levels",
	Long:  "manipulate YouTube memberships levels, only list for now",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)
}

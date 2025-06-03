package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "List memberships levels' info"
	long        = "List memberships levels' info, such as id, displayName, etc"
	partsUsage  = "Comma separated parts"
	outputUsage = "json or yaml"
)

var (
	parts  []string
	output string
)

var membershipsLevelCmd = &cobra.Command{
	Use:   "membershipsLevel",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)
}

package watermark

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	channelId              string
	file                   string
	inVideoPosition        string
	durationMs             uint64
	offsetMs               uint64
	offsetType             string
	onBehalfOfContentOwner string
)

var wartermarkCmd = &cobra.Command{
	Use:   "watermark",
	Short: "manipulate Youtube watermarks",
	Long:  "set or unset Youtube watermarks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(wartermarkCmd)
}

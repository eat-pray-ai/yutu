package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	file string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var channelBannerCmd = &cobra.Command{
	Use:   "channelBanner",
	Short: "insert Youtube channelBanner",
	Long:  "insert Youtube channelBanner",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelBannerCmd)
}

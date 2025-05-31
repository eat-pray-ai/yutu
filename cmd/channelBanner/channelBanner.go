package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	file   string
	output string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var channelBannerCmd = &cobra.Command{
	Use:   "channelBanner",
	Short: "Insert Youtube channel banner",
	Long:  "Insert Youtube channel banner",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelBannerCmd)
}

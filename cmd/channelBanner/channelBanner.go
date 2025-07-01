package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Insert Youtube channel banner"
	long      = "Insert Youtube channel banner"
	cidUsage  = "ID of the channel to insert the banner for"
	fileUsage = "Path to the banner image"
)

var (
	channelId string
	file      string
	output    string
	jpath     string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var channelBannerCmd = &cobra.Command{
	Use:   "channelBanner",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelBannerCmd)
}

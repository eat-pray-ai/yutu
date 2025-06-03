package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "Insert Youtube channel banner"
	long        = "Insert Youtube channel banner"
	fileUsage   = "Path to the banner image"
	outputUsage = "json, yaml, or silent"
)

var (
	file   string
	output string

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

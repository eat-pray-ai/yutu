package channelBanner

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/channelBanner"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a ChannelBanner",
	Long:  "Insert a ChannelBanner",
	Run: func(cmd *cobra.Command, args []string) {
		cb := channelBanner.NewChannelBanner(
			channelBanner.WithFile(file),
			channelBanner.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channelBanner.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		)
		cb.Insert()
	},
}

func init() {
	channelBannerCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the banner image")
	insertCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	insertCmd.Flags().StringVarP(&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "")

	insertCmd.MarkFlagRequired("file")
}

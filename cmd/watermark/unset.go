package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/watermark"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "unset watermark for channel's video",
	Long:  "unset watermark for channel's video",
	Run: func(cmd *cobra.Command, args []string) {
		w := watermark.NewWatermark(watermark.WithChannelId(channelId))
		w.Unset()
	},
}

func init() {
	wartermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "i", "", "ID of the channel to set watermark")
	unsetCmd.MarkFlagRequired("channelId")
}

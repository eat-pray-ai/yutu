package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/spf13/cobra"
)

const (
	unsetShort = "Unset watermark for channel's video"
	unsetLong  = "Unset watermark for channel's video by channel id"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: unsetShort,
	Long:  unsetLong,
	Run: func(cmd *cobra.Command, args []string) {
		w := watermark.NewWatermark(
			watermark.WithChannelId(channelId),
			watermark.WithService(nil),
		)
		w.Unset()
	},
}

func init() {
	watermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	_ = unsetCmd.MarkFlagRequired("channelId")
}

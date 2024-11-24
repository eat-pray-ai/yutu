package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset watermark for channel's video",
	Long:  "Unset watermark for channel's video",
	Run: func(cmd *cobra.Command, args []string) {
		w := watermark.NewWatermark(
			watermark.WithChannelId(channelId),
			watermark.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		w.Unset()
	},
}

func init() {
	wartermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of the channel to set watermark")
	unsetCmd.MarkFlagRequired("channelId")
}

package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set watermark for channel's video",
	Long:  "Set watermark for channel's video",
	Run: func(cmd *cobra.Command, args []string) {
		w := watermark.NewWatermark(
			watermark.WithChannelId(channelId),
			watermark.WithFile(file),
			watermark.WithInVideoPosition(inVideoPosition),
			watermark.WithDurationMs(durationMs),
			watermark.WithOffsetMs(offsetMs),
			watermark.WithOffsetType(offsetType),
			watermark.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			watermark.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		w.Set()
	},
}

func init() {
	wartermarkCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of channel to set watermark")
	setCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the watermark file")
	setCmd.Flags().StringVarP(
		&inVideoPosition, "inVideoPosition", "p", "", "topLeft, topRight, bottomLeft or bottomRight",
	)
	setCmd.Flags().Uint64VarP(
		&durationMs, "durationMs", "d", 0, "Duration in milliseconds for which the watermark should be displayed",
	)
	setCmd.Flags().Uint64VarP(&offsetMs, "offsetMs", "m", 0, "Defines the time at which the watermark will appear")
	setCmd.Flags().StringVarP(&offsetType, "offsetType", "t", "", "offsetFromStart or offsetFromEnd")
	setCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")

	setCmd.MarkFlagRequired("channelId")
	setCmd.MarkFlagRequired("file")
}

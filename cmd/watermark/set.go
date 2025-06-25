package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/spf13/cobra"
)

const (
	setShort = "Set watermark for channel's video"
	setLong  = "Set watermark for channel's video by channel id"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: setShort,
	Long:  setLong,
	Run: func(cmd *cobra.Command, args []string) {
		w := watermark.NewWatermark(
			watermark.WithChannelId(channelId),
			watermark.WithFile(file),
			watermark.WithInVideoPosition(inVideoPosition),
			watermark.WithDurationMs(durationMs),
			watermark.WithOffsetMs(offsetMs),
			watermark.WithOffsetType(offsetType),
			watermark.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			watermark.WithService(nil),
		)

		err := w.Set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	watermarkCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(
		&inVideoPosition, "inVideoPosition", "p", "", ivpUsage,
	)
	setCmd.Flags().Uint64VarP(&durationMs, "durationMs", "d", 0, dmUsage)
	setCmd.Flags().Uint64VarP(&offsetMs, "offsetMs", "m", 0, omUsage)
	setCmd.Flags().StringVarP(&offsetType, "offsetType", "t", "", otUsage)
	setCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = setCmd.MarkFlagRequired("channelId")
	_ = setCmd.MarkFlagRequired("file")
}

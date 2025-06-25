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

		err := w.Unset(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	watermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	_ = unsetCmd.MarkFlagRequired("channelId")
}

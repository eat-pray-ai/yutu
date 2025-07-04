package watermark

import (
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/spf13/cobra"
	"io"
)

const (
	unsetShort = "Unset watermark for channel's video"
	unsetLong  = "Unset watermark for channel's video by channel id"
)

func init() {
	watermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	_ = unsetCmd.MarkFlagRequired("channelId")
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: unsetShort,
	Long:  unsetLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := unset(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func unset(writer io.Writer) error {
	w := watermark.NewWatermark(
		watermark.WithChannelId(channelId),
		watermark.WithService(nil),
	)

	return w.Unset(writer)
}

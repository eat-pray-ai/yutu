package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/channelBanner"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	channelBannerCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("file")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insert(writer io.Writer) error {
	cb := channelBanner.NewChannelBanner(
		channelBanner.WithChannelId(channelId),
		channelBanner.WithFile(file),
		channelBanner.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channelBanner.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		channelBanner.WithService(nil),
	)

	return cb.Insert(output, jpath, writer)
}

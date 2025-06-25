package commentThread

import (
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/spf13/cobra"
)

const (
	insertShort    = "Insert a new comment thread"
	insertLong     = "Insert a new comment thread"
	insertVidUsage = "ID of the video"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		ct := commentThread.NewCommentThread(
			commentThread.WithAuthorChannelId(authorChannelId),
			commentThread.WithChannelId(channelId),
			commentThread.WithTextOriginal(textOriginal),
			commentThread.WithVideoId(videoId),
			commentThread.WithService(nil),
		)

		err := ct.Insert(output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", toUsage)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", insertVidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

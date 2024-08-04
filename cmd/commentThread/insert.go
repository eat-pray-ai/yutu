package commentThread

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/commentThread"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a new comment thread",
	Long:  "Insert a new comment thread",
	Run: func(cmd *cobra.Command, args []string) {
		ct := commentThread.NewCommentThread(
			commentThread.WithAuthorChannelId(authorChannelId),
			commentThread.WithChannelId(channelId),
			commentThread.WithTextOriginal(textOriginal),
			commentThread.WithVideoId(videoId),
			commentThread.WithService(nil),
		)
		ct.Insert(output)
	},
}

func init() {
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(&authorChannelId, "authorChannelId", "a", "", "Channel ID of the comment author")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "Channel ID of the video owner")
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", "Text of the comment")
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

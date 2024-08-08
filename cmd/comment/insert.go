package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a comment",
	Long:  "Insert a comment to a YouTube video",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithAuthorChannelId(authorChannelId),
			comment.WithChannelId(channelId),
			comment.WithCanRate(canRate, true),
			comment.WithParentId(parentId),
			comment.WithTextOriginal(textOriginal),
			comment.WithVideoId(videoId),
			comment.WithService(nil),
		)
		c.Insert(output)
	},
}

func init() {
	commentCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&authorChannelId, "authorChannelId", "a", "", "Channel ID of the comment author")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "Channel ID of the video owner")
	insertCmd.Flags().BoolVarP(&canRate, "canRate", "r", false, "Whether the viewer can rate the comment")
	insertCmd.Flags().StringVarP(&parentId, "parentId", "p", "", "ID of the parent comment")
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", "Text of the comment")
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

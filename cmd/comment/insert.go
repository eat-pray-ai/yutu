package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	insertShort       = "Insert a comment"
	insertLong        = "Insert a comment to a video"
	insertPidUsage    = "ID of the parent comment"
	insertOutputUsage = "json, yaml, or silent"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithAuthorChannelId(authorChannelId),
			comment.WithChannelId(channelId),
			comment.WithCanRate(canRate),
			comment.WithParentId(parentId),
			comment.WithTextOriginal(textOriginal),
			comment.WithVideoId(videoId),
			comment.WithService(nil),
		)

		err := c.Insert(output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	commentCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", cidUsage,
	)
	insertCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	insertCmd.Flags().StringVarP(
		&parentId, "parentId", "P", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", insertOutputUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("parentId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

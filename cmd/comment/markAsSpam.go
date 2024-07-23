package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/comment"
	"github.com/spf13/cobra"
)

var markAsSpamCmd = &cobra.Command{
	Use:   "markAsSpam",
	Short: "Mark YouTube comments as spam",
	Long:  "Mark YouTube comments as spam by ids",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithService(nil),
		)
		c.MarkAsSpam(false)
	},
}

func init() {
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, "Comma separated ids of comments")
}

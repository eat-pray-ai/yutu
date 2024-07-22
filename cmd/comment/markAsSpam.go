package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/comment"
	"github.com/spf13/cobra"
)

var markAsSpamCmd = &cobra.Command{
	Use:   "markAsSpam",
	Short: "Mark YouTube comments as spam",
	Long:  "Mark YouTube comments as spam by IDs",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(IDs),
			comment.WithService(nil),
		)
		c.MarkAsSpam(false)
	},
}

func init() {
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&IDs, "ids", "i", []string{}, "Comma separated IDs of comments")
}

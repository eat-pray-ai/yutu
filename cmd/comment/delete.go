package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/comment"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete YouTube comments",
	Long:  "Delete YouTube comments by IDs",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(IDs),
			comment.WithService(nil),
		)
		c.Delete()
	},
}

func init() {
	commentCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&IDs, "ids", "i", []string{}, "Comma separated IDs of comments")
}

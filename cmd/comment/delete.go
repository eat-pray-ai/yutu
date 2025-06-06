package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	deleteShort = "Delete YouTube comments"
	deleteLong  = "Delete YouTube comments by ids"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithService(nil),
		)
		c.Delete()
	},
}

func init() {
	commentCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

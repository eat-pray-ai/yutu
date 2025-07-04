package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort = "Delete YouTube comments"
	deleteLong  = "Delete YouTube comments by ids"
)

func init() {
	commentCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func del(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithService(nil),
	)

	return c.Delete(writer)
}

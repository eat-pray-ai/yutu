package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	masShort       = "Mark YouTube comments as spam"
	masLong        = "Mark YouTube comments as spam by ids"
	masOutputUsage = "json, yaml, or silent"
)

var markAsSpamCmd = &cobra.Command{
	Use:   "markAsSpam",
	Short: masShort,
	Long:  masLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithService(nil),
		)

		err := c.MarkAsSpam(output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	markAsSpamCmd.Flags().StringVarP(&output, "output", "o", "", masOutputUsage)

	_ = markAsSpamCmd.MarkFlagRequired("ids")
}

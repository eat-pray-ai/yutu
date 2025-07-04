package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
	"io"
)

const (
	masShort = "Mark YouTube comments as spam"
	masLong  = "Mark YouTube comments as spam by ids"
)

func init() {
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	markAsSpamCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	markAsSpamCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = markAsSpamCmd.MarkFlagRequired("ids")
}

var markAsSpamCmd = &cobra.Command{
	Use:   "markAsSpam",
	Short: masShort,
	Long:  masLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := markAsSpam(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func markAsSpam(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithService(nil),
	)

	return c.MarkAsSpam(output, jpath, writer)
}

package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	listShort       = "List YouTube comments"
	listLong        = "List YouTube comments by ids"
	listPidUsage    = "Returns replies to the specified comment"
	listOutputUsage = "json or yaml"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithMaxResults(maxResults),
			comment.WithParentId(parentId),
			comment.WithTextFormat(textFormat),
			comment.WithService(nil),
		)

		err := c.List(parts, output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	commentCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, mrUsage,
	)
	listCmd.Flags().StringVarP(&parentId, "parentId", "P", "", listPidUsage)
	listCmd.Flags().StringVarP(
		&textFormat, "textFormat", "t", "html", tfUsage,
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
}

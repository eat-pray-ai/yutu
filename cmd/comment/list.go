package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List YouTube comments",
	Long:  "List YouTube comments by ids",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithMaxResults(maxResults),
			comment.WithParentId(parentId),
			comment.WithTextFormat(textFormat),
			comment.WithService(nil),
		)
		c.List(parts, output)
	},
}

func init() {
	commentCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, "Comma separated ids of comments")
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(&parentId, "parentId", "P", "", "Returns replies to the specified comment")
	listCmd.Flags().StringVarP(&textFormat, "textFormat", "t", "", "textFormatUnspecified, html(default) or plainText")
	listCmd.Flags().StringSliceVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json or yaml")
}

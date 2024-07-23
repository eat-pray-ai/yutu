package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/comment"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a comment",
	Long:  "Update a comment on a YouTube video",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithCanRate(canRate, cmd.Flags().Lookup("canRate").Changed),
			comment.WithTextOriginal(textOriginal),
			comment.WithViewerRating(viewerRating),
			comment.WithService(nil),
		)
		c.Update(false)
	},
}

func init() {
	commentCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, "ID of the comment")
	updateCmd.Flags().BoolVarP(&canRate, "canRate", "r", false, "Whether the viewer can rate the comment")
	updateCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", "Text of the comment")
	updateCmd.Flags().StringVarP(&viewerRating, "viewerRating", "v", "", "none, like or dislike")
}

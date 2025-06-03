package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update a comment"
	updateLong        = "Update a comment on a YouTube video"
	updateIdUsage     = "ID of the comment"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithCanRate(canRate),
			comment.WithTextOriginal(textOriginal),
			comment.WithViewerRating(viewerRating),
			comment.WithService(nil),
		)
		c.Update(output)
	},
}

func init() {
	commentCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	updateCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	updateCmd.Flags().StringVarP(
		&viewerRating, "viewerRating", "r", "", vrUsage,
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)
}

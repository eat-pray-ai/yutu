package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

const (
	updateShort   = "Update a comment"
	updateLong    = "Update a comment on a video"
	updateIdUsage = "ID of the comment"
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

		err := c.Update(output, jpath, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
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
	updateCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

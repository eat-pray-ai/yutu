package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
)

var setModerationStatusCmd = &cobra.Command{
	Use:   "setModerationStatus",
	Short: "Set YouTube comments moderation status",
	Long:  "Set YouTube comments moderation status by ids",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(ids),
			comment.WithModerationStatus(moderationStatus),
			comment.WithBanAuthor(banAuthor, true),
		)
		c.SetModerationStatus(output)
	},
}

func init() {
	commentCmd.AddCommand(setModerationStatusCmd)

	setModerationStatusCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, "Comma separated ids of comments")
	setModerationStatusCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "s", "", "heldForReview, published or rejected",
	)
	setModerationStatusCmd.Flags().BoolVarP(&banAuthor, "banAuthor", "b", false, "true or false")
	setModerationStatusCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

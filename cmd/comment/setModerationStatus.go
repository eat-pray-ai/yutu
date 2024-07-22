package comment

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/comment"
	"github.com/spf13/cobra"
)

var setModerationStatusCmd = &cobra.Command{
	Use:   "setModerationStatus",
	Short: "Set YouTube comments moderation status",
	Long:  "Set YouTube comments moderation status by IDs",
	Run: func(cmd *cobra.Command, args []string) {
		c := comment.NewComment(
			comment.WithIDs(IDs),
			comment.WithModerationStatus(ModerationStatus),
			comment.WithBanAuthor(BanAuthor, true),
		)
		c.SetModerationStatus(false)
	},
}

func init() {
	commentCmd.AddCommand(setModerationStatusCmd)

	setModerationStatusCmd.Flags().StringSliceVarP(&IDs, "ids", "i", []string{}, "Comma separated IDs of comments")
	setModerationStatusCmd.Flags().StringVarP(
		&ModerationStatus, "moderationStatus", "s", "", "heldForReview, published or rejected",
	)
	setModerationStatusCmd.Flags().BoolVarP(&BanAuthor, "banAuthor", "b", false, "true or false")
}

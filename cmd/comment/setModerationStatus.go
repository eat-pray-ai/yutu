package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/spf13/cobra"
	"io"
)

const (
	smsShort = "Set YouTube comments moderation status"
	smsLong  = "Set YouTube comments moderation status by ids"
)

func init() {
	commentCmd.AddCommand(setModerationStatusCmd)

	setModerationStatusCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, idsUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "s", "", msUsage,
	)
	setModerationStatusCmd.Flags().BoolVarP(
		banAuthor, "banAuthor", "A", false, baUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&output, "output", "o", "", cmd.SilentUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&jpath, "jsonpath", "j", "", cmd.JpUsage,
	)

	_ = setModerationStatusCmd.MarkFlagRequired("ids")
	_ = setModerationStatusCmd.MarkFlagRequired("moderationStatus")
}

var setModerationStatusCmd = &cobra.Command{
	Use:   "setModerationStatus",
	Short: smsShort,
	Long:  smsLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := setModerationStatus(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setModerationStatus(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithModerationStatus(moderationStatus),
		comment.WithBanAuthor(banAuthor),
	)

	return c.SetModerationStatus(output, jpath, writer)
}

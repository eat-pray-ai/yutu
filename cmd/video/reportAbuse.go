package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
	"io"
)

const (
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video"
	raIdsUsage       = "IDs of the videos to report abuse on"
	raLangUsage      = "Language that the content was viewed in"
)

func init() {
	videoCmd.AddCommand(reportAbuseCmd)

	reportAbuseCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, raIdsUsage,
	)
	reportAbuseCmd.Flags().StringVarP(&reasonId, "reasonId", "r", "", ridUsage)
	reportAbuseCmd.Flags().StringVarP(
		&secondaryReasonId, "secondaryReasonId", "s", "", sridUsage,
	)
	reportAbuseCmd.Flags().StringVarP(
		&comments, "comments", "c", "", commentsUsage,
	)
	reportAbuseCmd.Flags().StringVarP(&language, "language", "l", "", raLangUsage)
	reportAbuseCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = reportAbuseCmd.MarkFlagRequired("ids")
	_ = reportAbuseCmd.MarkFlagRequired("reasonId")
}

var reportAbuseCmd = &cobra.Command{
	Use:   "reportAbuse",
	Short: reportAbuseShort,
	Long:  reportAbuseLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := reportAbuse(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func reportAbuse(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithReasonId(reasonId),
		video.WithSecondaryReasonId(secondaryReasonId),
		video.WithComments(comments),
		video.WithLanguage(language),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithService(nil),
	)

	return v.ReportAbuse(writer)
}

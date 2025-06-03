package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	reportAbuseShort = "Report abuse on a video"
	reportAbuseLong  = "Report abuse on a video"
	raIdUsage        = "ID of the video to report abuse on"
	raLangUsage      = "Language that the content was viewed in"
)

var reportAbuseCmd = &cobra.Command{
	Use:   "reportAbuse",
	Short: reportAbuseShort,
	Long:  reportAbuseLong,
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id),
			video.WithReasonId(reasonId),
			video.WithSecondaryReasonId(secondaryReasonId),
			video.WithComments(comments),
			video.WithLanguage(language),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithService(nil),
		)
		v.ReportAbuse()
	},
}

func init() {
	videoCmd.AddCommand(reportAbuseCmd)

	reportAbuseCmd.Flags().StringVarP(&id, "id", "i", "", raIdUsage)
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
}

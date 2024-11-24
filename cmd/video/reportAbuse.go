package video

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

var reportAbuseCmd = &cobra.Command{
	Use:   "reportAbuse",
	Short: "Report abuse on a video",
	Long:  "Report abuse on a video",
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithID(id),
			video.WithReasonId(reasonId),
			video.WithSecondaryReasonId(secondaryReasonId),
			video.WithComments(comments),
			video.WithLanguage(language),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		v.ReportAbuse()
	},
}

func init() {
	videoCmd.AddCommand(reportAbuseCmd)

	reportAbuseCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video to report abuse on")
	reportAbuseCmd.Flags().StringVarP(&reasonId, "reasonId", "r", "", "ID of the reason for reporting abuse")
	reportAbuseCmd.Flags().StringVarP(
		&secondaryReasonId, "secondaryReasonId", "s", "", "ID of the secondary reason for reporting abuse",
	)
	reportAbuseCmd.Flags().StringVarP(&comments, "comments", "c", "", "Additional comments regarding the abuse report")
	reportAbuseCmd.Flags().StringVarP(&language, "language", "l", "", "The language that the content was viewed in")
	reportAbuseCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
}

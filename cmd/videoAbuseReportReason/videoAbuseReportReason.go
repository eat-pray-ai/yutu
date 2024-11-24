package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	hl         string
	parts      []string
	output     string
	credential string
	cacheToken string
)

var videoAbuseReportReasonCmd = &cobra.Command{
	Use:   "videoAbuseReportReason",
	Short: "List YouTube video abuse report reasons",
	Long:  "List YouTube video abuse report reasons",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoAbuseReportReasonCmd)

	videoAbuseReportReasonCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	videoAbuseReportReasonCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "List YouTube video abuse report reasons"
	long        = "List YouTube video abuse report reasons"
	hlUsage     = "Host language"
	partsUsage  = "Comma separated parts"
	outputUsage = "json or yaml"
)

var (
	hl     string
	parts  []string
	output string
)

var videoAbuseReportReasonCmd = &cobra.Command{
	Use:   "videoAbuseReportReason",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoAbuseReportReasonCmd)
}

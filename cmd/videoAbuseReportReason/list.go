package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/pkg/videoAbuseReportReason"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		va := videoAbuseReportReason.NewVideoAbuseReportReason(
			videoAbuseReportReason.WithHL(hl),
			videoAbuseReportReason.WithService(nil),
		)

		err := va.List(parts, output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	videoAbuseReportReasonCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", outputUsage)
}

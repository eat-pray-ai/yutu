package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list YouTube video abuse report reasons",
	Long:  "list YouTube video abuse report reasons",
	Run: func(cmd *cobra.Command, args []string) {
		va := yutuber.NewVideoAbuseReportReason(
			yutuber.WithVideoAbuseReportReasonHL(hl),
		)
		va.List(parts, output)
	},
}

func init() {
	videoAbuseReportReasonCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Host language",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/videoAbuseReportReason"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	videoAbuseReportReasonCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func list(writer io.Writer) error {
	va := videoAbuseReportReason.NewVideoAbuseReportReason(
		videoAbuseReportReason.WithHL(hl),
		videoAbuseReportReason.WithService(nil),
	)

	return va.List(parts, output, jpath, writer)
}

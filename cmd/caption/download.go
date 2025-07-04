package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
	"io"
)

const (
	downloadShort   = "Download caption"
	downloadLong    = "Download caption from a video"
	downloadIdUsage = "ID of the caption to download"
)

func init() {
	captionCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringSliceVarP(
		&ids, "id", "i", []string{}, downloadIdUsage,
	)
	downloadCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", tfmtUsage)
	downloadCmd.Flags().StringVarP(&tlang, "tlang", "l", "", tlangUsage)
	downloadCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	_ = downloadCmd.MarkFlagRequired("id")
	_ = downloadCmd.MarkFlagRequired("file")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: downloadShort,
	Long:  downloadLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := download(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func download(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithIDs(ids),
		caption.WithFile(file),
		caption.WithTfmt(tfmt),
		caption.WithTlang(tlang),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithService(nil),
	)

	return c.Download(writer)
}

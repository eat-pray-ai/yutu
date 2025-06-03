package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

const (
	downloadShort = "Download caption"
	downloadLong  = "Download caption from a video"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: downloadShort,
	Long:  downloadLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithID(id),
			caption.WithFile(file),
			caption.WithTfmt(tfmt),
			caption.WithTlang(tlang),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithService(nil),
		)
		c.Download()
	},
}

func init() {
	captionCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&id, "id", "i", "", idUsage)
	downloadCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", tfmtUsage)
	downloadCmd.Flags().StringVarP(&tlang, "tlang", "l", "", tlangUsage)
	downloadCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	downloadCmd.MarkFlagRequired("id")
	downloadCmd.MarkFlagRequired("file")
}

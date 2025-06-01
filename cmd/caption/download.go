package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download caption",
	Long:  "Download caption from a video",
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

	downloadCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the caption")
	downloadCmd.Flags().StringVarP(
		&file, "file", "f", "", "Path to save the caption file",
	)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", "sbv, srt or vtt")
	downloadCmd.Flags().StringVarP(
		&tlang, "tlang", "l", "", "Translate the captions into this language",
	)
	downloadCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	downloadCmd.MarkFlagRequired("id")
	downloadCmd.MarkFlagRequired("file")
}

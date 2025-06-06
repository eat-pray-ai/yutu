package thumbnail

import (
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		t := thumbnail.NewThumbnail(
			thumbnail.WithFile(file),
			thumbnail.WithVideoId(videoId),
			thumbnail.WithService(nil),
		)
		t.Set(output)
	},
}

func init() {
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)

	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

package thumbnail

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/thumbnail"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set thumbnail for a video",
	Long:  "Set thumbnail for a video",
	Run: func(cmd *cobra.Command, args []string) {
		t := thumbnail.NewThumbnail(
			thumbnail.WithFile(file),
			thumbnail.WithVideoId(videoId),
			thumbnail.WithService(nil),
		)
		t.Set(false)
	},
}

func init() {
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringP("file", "f", "", "Path to the thumbnail file")
	setCmd.Flags().StringP("videoId", "v", "", "ID of the video")
}

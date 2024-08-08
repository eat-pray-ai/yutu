package thumbnail

import (
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
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
		t.Set(output)
	},
}

func init() {
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the thumbnail file")
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	setCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

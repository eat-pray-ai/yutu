package playlistImage

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/spf13/cobra"
)

const (
	insertShort       = "Insert a YouTube playlist image"
	insertLong        = "Insert a YouTube playlist image for a given playlist id"
	insertOutputUsage = "json, yaml, or silent"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithFile(file),
			playlistImage.WithPlaylistID(playlistId),
			playlistImage.WithType(type_),
			playlistImage.WithHeight(height),
			playlistImage.WithWidth(width),
			playlistImage.WithService(nil),
		)
		pi.Insert(output)
	},
}

func init() {
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", insertOutputUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

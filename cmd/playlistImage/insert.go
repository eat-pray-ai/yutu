package playlistImage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort = "Insert a YouTube playlist image"
	insertLong  = "Insert a YouTube playlist image for a given playlist id"
)

func init() {
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insert(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithFile(file),
		playlistImage.WithPlaylistID(playlistId),
		playlistImage.WithType(type_),
		playlistImage.WithHeight(height),
		playlistImage.WithWidth(width),
		playlistImage.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

package thumbnail

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Set thumbnail for a video"
	long      = "Set thumbnail for a video"
	fileUsage = "Path to the thumbnail file"
	vidUsage  = "ID of the video"
)

var (
	file    string
	videoId string
	output  string
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(thumbnailCmd)
}

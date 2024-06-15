package thumbnail

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	file    string
	videoId string
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: "Set thumbnail for a video",
	Long:  "Set thumbnail for a video",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(thumbnailCmd)
}

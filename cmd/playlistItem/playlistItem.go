package playlistItem

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id         string
	title      string
	desc       string
	videoId    string
	playlistId string
	channelId  string
	privacy    string
	output     string
	parts      []string
)

var playlistItemCmd = &cobra.Command{
	Use:   "playlistItem",
	Short: "manipulate YouTube playlist items",
	Long:  "manipulate YouTube playlist items, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistItemCmd)
}

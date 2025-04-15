package playlistItem

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id          string
	title       string
	description string
	kind        string
	kVideoId    string
	kChannelId  string
	kPlaylistId string
	videoId     string
	playlistId  string
	channelId   string
	maxResults  int64
	privacy     string
	output      string
	parts       []string

	onBehalfOfContentOwner string
)

var playlistItemCmd = &cobra.Command{
	Use:   "playlistItem",
	Short: "Manipulate YouTube playlist items",
	Long:  "List, insert, update, or delete YouTube playlist items",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistItemCmd)
}

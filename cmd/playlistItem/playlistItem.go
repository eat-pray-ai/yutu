package playlistItem

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short        = "Manipulate YouTube playlist items"
	long         = "List, insert, update, or delete YouTube playlist items"
	titleUsage   = "Title of the playlist item"
	descUsage    = "Description of the playlist item"
	kindUsage    = "video, channel, or playlist"
	kvidUsage    = "ID of the video if kind is video"
	kcidUsage    = "ID of the channel if kind is channel"
	kpidUsage    = "ID of the playlist if kind is playlist"
	vidUsage     = "Return the playlist items associated with the given video id"
	cidUsage     = "ID that YouTube uses to uniquely identify the user that added the item to the playlist"
	privacyUsage = "public, private, or unlisted"
)

var (
	ids         []string
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
	parts       []string
	output      string
	jpath       string

	onBehalfOfContentOwner string
)

var playlistItemCmd = &cobra.Command{
	Use:   "playlistItem",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistItemCmd)
}

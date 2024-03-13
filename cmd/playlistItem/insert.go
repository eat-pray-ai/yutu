package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "insert playlist item",
	Long:  "insert playlist item into a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		pi := yutuber.NewPlaylistItem(
			yutuber.WithPlaylistItemTitle(title),
			yutuber.WithPlaylistItemDesc(desc),
			yutuber.WithPlaylistItemVideoId(videoId),
			yutuber.WithPlaylistItemPlaylistId(playlistId),
			yutuber.WithPlaylistItemPrivacy(privacy),
		)
		pi.Insert()
	},
}

func init() {
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist item")
	insertCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the playlist item")
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", "ID of the playlist")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of the channel")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "r", "", "Privacy status of the playlist item")
}

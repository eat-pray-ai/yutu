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
			yutuber.WithPlaylistItemDescription(description),
			yutuber.WithPlaylistItemKind(kind),
			yutuber.WithPlaylistItemKVideoId(kVideoId),
			yutuber.WithPlaylistItemKChannelId(kChannelId),
			yutuber.WithPlaylistItemKPlaylistId(kPlaylistId),
			yutuber.WithPlaylistItemPlaylistId(playlistId),
			yutuber.WithPlaylistItemPrivacy(privacy),
		)
		pi.Insert()
	},
}

func init() {
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&title, "title", "t", "", "Title of the playlist item",
	)
	insertCmd.Flags().StringVarP(
		&description, "description", "d", "", "Description of the playlist item",
	)
	insertCmd.Flags().StringVarP(
		&kind, "kind", "k", "", "video, channel, or playlist",
	)
	insertCmd.Flags().StringVarP(
		&kVideoId, "kVideoId", "V", "", "ID of the video if kind is video",
	)
	insertCmd.Flags().StringVarP(
		&kChannelId, "kChannelId", "C", "", "ID of the channel if kind is channel",
	)
	insertCmd.Flags().StringVarP(
		&kPlaylistId, "kPlaylistId", "P", "",
		"ID of the playlist if kind is playlist",
	)
	insertCmd.Flags().StringVarP(
		&playlistId, "playlistId", "p", "",
		"The ID that YouTube uses to uniquely identify the playlist that the playlist item is in",
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "",
		"The ID that YouTube uses to uniquely identify the user that added the item to the playlist",
	)
	insertCmd.Flags().StringVarP(
		&privacy, "privacy", "r", "", "public, private, or unlisted",
	)
}

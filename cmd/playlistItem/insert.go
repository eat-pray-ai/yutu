package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert a playlist item into a playlist",
	Long:  "Insert a playlist item into a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(title),
			playlistItem.WithDescription(description),
			playlistItem.WithKind(kind),
			playlistItem.WithKVideoId(kVideoId),
			playlistItem.WithKChannelId(kChannelId),
			playlistItem.WithKPlaylistId(kPlaylistId),
			playlistItem.WithPlaylistId(playlistId),
			playlistItem.WithPrivacy(privacy),
			playlistItem.WithChannelId(channelId),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithService(auth.NewY2BService(
				auth.WithCredential(credential),
				auth.WithCacheToken(cacheToken),
			)),
		)
		pi.Insert(output)
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
		&kPlaylistId, "kPlaylistId", "Y", "",
		"ID of the playlist if kind is playlist",
	)
	insertCmd.Flags().StringVarP(
		&playlistId, "playlistId", "y", "",
		"The ID that YouTube uses to uniquely identify the playlist that the playlist item is in",
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "",
		"The ID that YouTube uses to uniquely identify the user that added the item to the playlist",
	)
	insertCmd.Flags().StringVarP(
		&privacy, "privacy", "p", "", "public, private, or unlisted",
	)
	insertCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

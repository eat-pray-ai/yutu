package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "create a new playlist",
	Long:  "create a new playlist, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := yutuber.NewPlaylist(
			yutuber.WithPlaylistTitle(title),
			yutuber.WithPlaylistDesc(desc),
			yutuber.WithPlaylistTags(tags),
			yutuber.WithPlaylistLanguage(language),
			yutuber.WithPlaylistChannelId(channelId),
			yutuber.WithPlaylistPrivacy(privacy),
		)
		p.Insert()
	},
}

func init() {
	playlistCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist")
	insertCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the playlist")
	insertCmd.Flags().StringVarP(&tags, "tags", "g", "", "Comma separated tags")
	insertCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the playlist")
	insertCmd.Flags().StringVarP(&channelId, "channel", "c", "", "Channel ID of the playlist")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "Privacy status of the playlist")

	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("channel")
	insertCmd.MarkFlagRequired("privacy")
}

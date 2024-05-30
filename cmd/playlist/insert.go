package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlist"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "create a new playlist",
	Long:  "create a new playlist, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithTitle(title),
			playlist.WithDescription(description),
			playlist.WithTags(tags),
			playlist.WithLanguage(language),
			playlist.WithChannelId(channelId),
			playlist.WithPrivacy(privacy),
		)
		p.Insert()
	},
}

func init() {
	playlistCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist")
	insertCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the playlist")
	insertCmd.Flags().StringArrayVarP(&tags, "tags", "a", []string{}, "Comma separated tags")
	insertCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the playlist")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "Channel ID of the playlist")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "public, private or unlisted")

	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("channel")
	insertCmd.MarkFlagRequired("privacy")
}

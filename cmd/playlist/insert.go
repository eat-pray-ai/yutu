package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
)

const (
	insertShort       = "Create a new playlist"
	insertLong        = "Create a new playlist, with the specified title, description, tags, etc"
	insertCidUsage    = "Channel id of the playlist"
	insertOutputUsage = "json, yaml, or silent"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithTitle(title),
			playlist.WithDescription(description),
			playlist.WithTags(tags),
			playlist.WithLanguage(language),
			playlist.WithChannelId(channelId),
			playlist.WithPrivacy(privacy),
			playlist.WithService(nil),
		)
		p.Insert(output)
	},
}

func init() {
	playlistCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringArrayVarP(&tags, "tags", "a", []string{}, tagsUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", insertOutputUsage)

	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("channel")
	insertCmd.MarkFlagRequired("privacy")
}

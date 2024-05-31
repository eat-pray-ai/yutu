package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlist"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an existing playlist",
	Long:  "update an existing playlist, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithId(id),
			playlist.WithTitle(title),
			playlist.WithDescription(description),
			playlist.WithTags(tags),
			playlist.WithLanguage(language),
			playlist.WithPrivacy(privacy),
			playlist.WithService(),
		)
		p.Update()
	},
}

func init() {
	playlistCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the playlist")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist")
	updateCmd.Flags().StringVarP(
		&description, "description", "d", "", "Description of the playlist",
	)
	updateCmd.Flags().StringArrayVarP(
		&tags, "tags", "a", []string{}, "Comma separated tags",
	)
	updateCmd.Flags().StringVarP(
		&language, "language", "l", "", "Language of the playlist",
	)
	updateCmd.Flags().StringVarP(
		&privacy, "privacy", "p", "", "Privacy status of the playlist",
	)

	updateCmd.MarkFlagRequired("id")
}

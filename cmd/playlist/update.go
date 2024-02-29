package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an existing playlist",
	Long:  "update an existing playlist, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := yutuber.NewPlaylist(
			yutuber.WithPlaylistId(id),
			yutuber.WithPlaylistTitle(title),
			yutuber.WithPlaylistDesc(desc),
			yutuber.WithPlaylistTags(tags),
			yutuber.WithPlaylistLanguage(language),
			yutuber.WithPlaylistPrivacy(privacy),
		)
		p.Update()
	},
}

func init() {
	playlistCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the playlist")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist")
	updateCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the playlist")
	updateCmd.Flags().StringArrayVarP(&tags, "tags", "a", []string{}, "Comma separated tags")
	updateCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the playlist")
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "Privacy status of the playlist")

	updateCmd.MarkFlagRequired("id")
}

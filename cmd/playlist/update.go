package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update an existing playlist"
	updateLong        = "Update an existing playlist, with the specified title, description, tags, etc"
	updateIdUsage     = "ID of the playlist to update"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		p := playlist.NewPlaylist(
			playlist.WithIDs(ids),
			playlist.WithTitle(title),
			playlist.WithDescription(description),
			playlist.WithTags(tags),
			playlist.WithLanguage(language),
			playlist.WithPrivacy(privacy),
			playlist.WithService(nil),
		)
		p.Update(output)
	},
}

func init() {
	playlistCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

package playlistItem

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/playlistItem"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update playlist item",
	Long:  "update playlist item's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistItem.NewPlaylistItem(
			playlistItem.WithId(id),
			playlistItem.WithTitle(title),
			playlistItem.WithDescription(description),
			playlistItem.WithPrivacy(privacy),
		)
		pi.Update()
	},
}

func init() {
	playlistItemCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the playlist item")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the playlist item")
	updateCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the playlist item")
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "Privacy status of the playlist item")
}

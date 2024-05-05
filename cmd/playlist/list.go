package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list playlist's info",
	Long:  "list playlist's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		p := yutuber.NewPlaylist(
			yutuber.WithPlaylistId(id),
			yutuber.WithPlaylistChannelId(channelId),
		)
		p.List(parts, output)
	},
}

func init() {
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "Return the playlists with the given IDs for Stubby or Apiary.")
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "Return the playlists owned by the specified channel ID")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet", "status"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
}

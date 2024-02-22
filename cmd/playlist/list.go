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
		p.List()
	},
}

func init() {
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the playlist")
	listCmd.Flags().StringVarP(&channelId, "channel", "c", "", "ID of the channel")
}

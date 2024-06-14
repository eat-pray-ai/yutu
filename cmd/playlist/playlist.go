package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

var (
	id          string
	title       string
	description string
	hl          string
	maxResults  int64
	mine        bool
	tags        []string
	language    string
	channelId   string
	privacy     string
	output      string
	parts       []string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "manipulate YouTube playlists",
	Long:  "manipulate YouTube playlists, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}

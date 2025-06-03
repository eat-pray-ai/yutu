package playlistImage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

// todo
var (
	id         string
	kind       string
	height     int64
	playlistId string
	type_      string
	width      int64
	file       string
	parent     string
	maxResults int64
	parts      []string
	output     string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var playlistImageCmd = &cobra.Command{
	Use:   "playlistImage",
	Short: "Manipulate YouTube playlist images",
	Long:  "List, insert, update, or delete YouTube playlist images",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistImageCmd)
}

package playlist

import (
	"fmt"

	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

var (
	id       string
	title    string
	desc     string
	tags     string
	language string
	channel  string
	privacy  string
	output   string
	parts    []string
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "manipulate YouTube playlists",
	Long:  "manipulate YouTube playlists, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("playlist called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}

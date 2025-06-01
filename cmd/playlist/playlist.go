package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	id          string
	title       string
	description string
	hl          string
	maxResults  int64
	mine        = utils.BoolPtr("false")
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
	Short: "Manipulate YouTube playlists",
	Long:  "List, insert, update, or delete YouTube playlists",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ResetBool(map[string]*bool{"mine": mine}, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}

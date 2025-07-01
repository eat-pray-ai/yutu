package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	short         = "Manipulate YouTube playlists"
	long          = "List, insert, update, or delete YouTube playlists"
	titleUsage    = "Title of the playlist"
	descUsage     = "Description of the playlist"
	hlUsage       = "Return content in specified language"
	mrUsage       = "The maximum number of items that should be returned"
	mineUsage     = "Return the playlists owned by the authenticated user"
	tagsUsage     = "Comma separated tags"
	languageUsage = "Language of the playlist"
	privacyUsage  = "public, private, or unlisted"
	partsUsage    = "Comma separated parts"
)

var (
	ids         []string
	title       string
	description string
	hl          string
	maxResults  int64
	mine        = utils.BoolPtr("false")
	tags        []string
	language    string
	channelId   string
	privacy     string
	parts       []string
	output      string
	jpath       string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ResetBool(map[string]**bool{"mine": &mine}, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}

package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	ids              []string
	authorChannelId  string
	canRate          = utils.BoolPtr("false")
	channelId        string
	maxResults       int64
	parentId         string
	textFormat       string
	textOriginal     string
	moderationStatus string
	banAuthor        = utils.BoolPtr("false")
	videoId          string
	viewerRating     string
	parts            []string
	output           string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manipulate YouTube comments",
	Long:  "List, insert, update, mark as spam, set moderation status, or delete YouTube comments",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]*bool{"canRate": canRate, "banAuthor": banAuthor}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentCmd)
}

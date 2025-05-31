package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	ids              []string
	authorChannelId  string
	canRate          bool
	channelId        string
	maxResults       int64
	parentId         string
	textFormat       string
	textOriginal     string
	moderationStatus string
	banAuthor        bool
	videoId          string
	viewerRating     string
	parts            []string
	output           string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manipulate YouTube comments",
	Long:  "List, insert, update, mark as spam, set moderation status, or delete YouTube comments",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentCmd)
}

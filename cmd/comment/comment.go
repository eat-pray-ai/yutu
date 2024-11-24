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
	credential       string
	cacheToken       string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manipulate YouTube comments",
	Long:  "Manipulate YouTube comments, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentCmd)

	commentCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	commentCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

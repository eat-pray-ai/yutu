package commentThread

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id                           []string
	allThreadsRelatedToChannelId string
	authorChannelId              string
	channelId                    string
	maxResults                   int64
	moderationStatus             string
	order                        string
	searchTerms                  string
	textFormat                   string
	textOriginal                 string
	videoId                      string
	parts                        []string
	output                       string
)

var commentThreadCmd = &cobra.Command{
	Use:   "commentThread",
	Short: "Manipulate YouTube comment threads",
	Long:  "List or insert YouTube comment threads",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentThreadCmd)
}

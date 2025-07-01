package commentThread

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short        = "Manipulate YouTube comment threads"
	long         = "List or insert YouTube comment threads"
	idsUsage     = "IDs of the comment threads"
	atrtcidUsage = "Returns the comment threads of all videos of the channel and the channel comments as well"
	acidUsage    = "Channel id of the comment author"
	cidUsage     = "Channel id of the video owner"
	mrUsage      = "The maximum number of items that should be returned"
	msUsage      = "published, heldForReview, likelySpam, or rejected"
	orderUsage   = "orderUnspecified, time, or relevance"
	stUsage      = "Search terms"
	tfUsage      = "textFormatUnspecified or html"
	toUsage      = "Text of the comment"
	partsUsage   = "Comma separated parts"
)

var (
	ids                          []string
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
	jpath                        string
)

var commentThreadCmd = &cobra.Command{
	Use:   "commentThread",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentThreadCmd)
}

package video

import (
	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

var (
	id                string
	autoLevels        bool
	file              string
	title             string
	description       string
	hl                string
	tags              []string
	language          string
	locale            string
	license           string
	thumbnail         string
	rating            string
	chart             string
	channelId         string
	comments          string
	playListId        string
	categoryId        string
	privacy           string
	forKids           bool
	embeddable        bool
	output            string
	parts             []string
	publishAt         string
	regionCode        string
	reasonId          string
	secondaryReasonId string
	stabilize         bool
	maxHeight         int64
	maxWidth          int64
	maxResults        int64

	notifySubscribers             bool
	publicStatsViewable           bool
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	credential                    string
	cacheToken                    string
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Manipulate YouTube videos",
	Long:  "Manipulate YouTube videos, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCmd)

	videoCmd.PersistentFlags().StringVarP(&credential, "credential", "c", "client_secret.json", "Path to client secret file")
	videoCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file")
}

package search

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           bool
	forDeveloper              bool
	forMine                   bool
	location                  string
	locationRadius            string
	maxResults                int64
	onBehalfOfContentOwner    string
	order                     string
	publishedAfter            string
	publishedBefore           string
	q                         string
	regionCode                string
	relevanceLanguage         string
	safeSearch                string
	topicId                   string
	types                     string
	videoCaption              string
	videoCategoryId           string
	videoDefinition           string
	videoDimension            string
	videoDuration             string
	videoEmbeddable           string
	videoLicense              string
	videoPaidProductPlacement string
	videoSyndicated           string
	videoType                 string
	parts                     []string
	output                    string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for YouTube resources",
	Long:  "Search for YouTube resources",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(searchCmd)
}

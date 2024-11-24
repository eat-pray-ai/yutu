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
	credential                string
	cacheToken                string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for Youtube resources",
	Long:  "Search for Youtube resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().StringVarP(&credential, "credential", "c", "client_secret.json", "Path to client secret file")
	searchCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file")
}

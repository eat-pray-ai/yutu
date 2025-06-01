package search

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           = utils.BoolPtr("")
	forDeveloper              = utils.BoolPtr("")
	forMine                   = utils.BoolPtr("")
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]*bool{
			"forContentOwner": forContentOwner,
			"forDeveloper":    forDeveloper,
			"forMine":         forMine,
		}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(searchCmd)
}

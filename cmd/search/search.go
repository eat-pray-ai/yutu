package search

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short         = "Search for YouTube resources"
	long          = "Search for YouTube resources"
	cidUsage      = "Filter on resources belonging to this channelId"
	ctUsage       = "channelTypeUnspecified, any, or show"
	etUsage       = "none, upcoming, live, or completed"
	fcoUsage      = "Search owned by content owner"
	fdUsage       = "Only retrieve videos uploaded using the project id of the authenticated user"
	fmUsage       = "Search for the private videos of the authenticated user"
	locationUsage = "Filter on location of the video"
	lrUsage       = "Filter on distance from the location"
	mrUsage       = "The maximum number of items that should be returned"
	orderUsage    = "searchSortUnspecified, date, rating, viewCount, relevance, title, videoCount"
	paUsage       = "Filter on resources published after this date"
	pbUsage       = "Filter on resources published before this date"
	qUsage        = "Textual search terms to match"
	rcUsage       = "Display the content as seen by viewers in this country"
	rlUsage       = "Return results relevant to this language"
	ssUsage       = "safeSearchSettingUnspecified, none, moderate, or strict"
	tidUsage      = "Restrict results to a particular topic"
	typesUsage    = "Restrict results to a particular set of resource types from One Platform"
	vcUsage       = "videoCaptionUnspecified, any, closedCaption, or none"
	vcidUsage     = "Filter on videos in a specific category"
	vdeUsage      = "Filter on the definition of the videos"
	vdiUsage      = "any, 2d, or 3d"
	vduUsage      = "videoDurationUnspecified, any, short, medium, or long"
	veUsage       = "videoEmbeddableUnspecified, any, or true"
	vlUsage       = "any, youtube, or creativeCommon"
	vpppUsage     = "videoPaidProductPlacementUnspecified, any, or true"
	vsUsage       = "videoSyndicatedUnspecified, any, or true"
	vtUsage       = "videoTypeUnspecified, any, movie, or episode"
	partsUsage    = "Comma separated parts"
	outputUsage   = "json, yaml, or table"
)

var (
	channelId                 string
	channelType               string
	eventType                 string
	forContentOwner           = utils.BoolPtr("false")
	forDeveloper              = utils.BoolPtr("false")
	forMine                   = utils.BoolPtr("false")
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
	types                     []string
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
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]**bool{
			"forContentOwner": &forContentOwner,
			"forDeveloper":    &forDeveloper,
			"forMine":         &forMine,
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

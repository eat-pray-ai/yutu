package video

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

var (
	id                string
	autoLevels        = utils.BoolPtr("")
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
	forKids           = utils.BoolPtr("")
	embeddable        = utils.BoolPtr("")
	output            string
	parts             []string
	publishAt         string
	regionCode        string
	reasonId          string
	secondaryReasonId string
	stabilize         = utils.BoolPtr("")
	maxHeight         int64
	maxWidth          int64
	maxResults        int64

	notifySubscribers             = utils.BoolPtr("")
	publicStatsViewable           = utils.BoolPtr("")
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Manipulate YouTube videos",
	Long:  "List, insert, update, rate, get rating, report abuse, or delete YouTube videos",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resetFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCmd)
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]*bool{
		"autoLevels":          autoLevels,
		"forKids":             forKids,
		"embeddable":          embeddable,
		"stabilize":           stabilize,
		"notifySubscribers":   notifySubscribers,
		"publicStatsViewable": publicStatsViewable,
	}

	utils.ResetBool(boolMap, flagSet)
}

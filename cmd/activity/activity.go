package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	shortUsage           = "List YouTube activities"
	longUsage            = "List YouTube activities"
	channelIdUsage       = "ID of the channel"
	homeUsage            = "true or false or empty"
	maxResultsUsage      = "The maximum number of items that should be returned"
	mineUsage            = "true or false or empty"
	publishedAfterUsage  = "Filter on activities published after this date"
	publishedBeforeUsage = "Filter on activities published before this date"
	regionCodeUsage      = ""
	partsUsage           = "Comma separated parts"
	outputUsage          = "json or yaml"
)

var (
	channelId       string
	home            = utils.BoolPtr("false")
	maxResults      int64
	mine            = utils.BoolPtr("false")
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
	output          string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: shortUsage,
	Long:  longUsage,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]*bool{
			"home": home,
			"mine": mine,
		}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

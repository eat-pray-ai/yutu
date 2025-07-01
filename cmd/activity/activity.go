package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short      = "List YouTube activities"
	long       = "List YouTube activities, such as likes, favorites, uploads, etc"
	ciUsage    = "ID of the channel"
	homeUsage  = "true, false, or empty"
	mrUsage    = "The maximum number of items that should be returned"
	mineUsage  = "true, false, or empty"
	paUsage    = "Filter on activities published after this date"
	pbUsage    = "Filter on activities published before this date"
	rcUsage    = ""
	partsUsage = "Comma separated parts"
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
	jpath           string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]**bool{
			"home": &home,
			"mine": &mine,
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

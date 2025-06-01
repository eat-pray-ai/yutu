package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	channelId       string
	home            = utils.BoolPtr("false")
	maxResults      int64
	mine            = utils.BoolPtr("true")
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
	output          string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "List YouTube activities",
	Long:  "List YouTube activities",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

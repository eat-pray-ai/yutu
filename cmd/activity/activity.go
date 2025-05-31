package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	channelId       string
	home            bool
	maxResults      int64
	mine            bool
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

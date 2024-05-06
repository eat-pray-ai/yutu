package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	channelId       string
	home            string
	maxResults      int64
	mine            string
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
	output          string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "manipulate YouTube activities",
	Long:  "manipulate YouTube activities, only list for now",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

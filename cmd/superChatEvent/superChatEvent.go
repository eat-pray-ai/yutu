package superChatEvent

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	hl         string
	maxResults int64
	parts      []string
	output     string
)

var superChatEventCmd = &cobra.Command{
	Use:   "superChatEvent",
	Short: "List Super Chat events for a YouTube channel",
	Long:  "List Super Chat events for a YouTube channel",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(superChatEventCmd)
}

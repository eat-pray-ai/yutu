package superChatEvent

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short      = "List Super Chat events for a channel"
	long       = "List Super Chat events for a channel"
	hlUsage    = "Return rendered funding amounts in specified language"
	mrUsage    = "The maximum number of items that should be returned"
	partsUsage = "Comma separated parts"
)

var (
	hl         string
	maxResults int64
	parts      []string
	output     string
	jpath      string
)

var superChatEventCmd = &cobra.Command{
	Use:   "superChatEvent",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(superChatEventCmd)
}

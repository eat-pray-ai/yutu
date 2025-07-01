package member

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short      = "List channel's members' info"
	long       = "List channel's members' info, such as channelId, displayName, etc"
	mcidUsage  = "Comma separated list of channel IDs. Only data about members that are part of this list will be included"
	hatlUsage  = "Filter members in the results set to the ones that have access to a level"
	mrUsage    = "The maximum number of items that should be returned"
	mmUsage    = "listMembersModeUnknown, updates, or all_current"
	partsUsage = "Comma separated parts"
)

var (
	memberChannelId  string
	hasAccessToLevel string
	maxResults       int64
	mode             string
	parts            []string
	output           string
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(memberCmd)
}

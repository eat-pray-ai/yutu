package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id                            string
	channelId                     string
	forChannelId                  string
	maxResult                     int64
	mine                          string
	myRecentSubscribers           string
	mySubscribers                 string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	parts                         []string
	output                        string
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "manipulate YouTube subscriptions",
	Long:  "manipulate YouTube subscriptions, such as list, insert and delete",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(subscriptionCmd)
}

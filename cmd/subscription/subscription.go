package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id                            string
	subscriberChannelId           string
	description                   string
	channelId                     string
	forChannelId                  string
	maxResults                    int64
	mine                          bool
	myRecentSubscribers           bool
	mySubscribers                 bool
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	title                         string
	parts                         []string
	output                        string
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Manipulate YouTube subscriptions",
	Long:  "List, insert, or delete YouTube subscriptions",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(subscriptionCmd)
}

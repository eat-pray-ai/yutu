// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short      = "Manipulate YouTube subscriptions"
	long       = "List, insert, or delete YouTube subscriptions"
	scidUsage  = "Subscriber's channel id"
	descUsage  = "Description of the subscription"
	fcidUsage  = "Return the subscriptions to the subset of these channels that the authenticated user is subscribed to"
	mineUsage  = "Return the subscriptions of the authenticated user"
	mrsUsage   = "true, false, or leave empty"
	msUsage    = "Return the subscribers of the given channel owner"
	orderUsage = "subscriptionOrderUnspecified, relevance, unread, or alphabetical"
	titleUsage = "Title of the subscription"
)

var (
	ids                           []string
	subscriberChannelId           string
	description                   string
	channelId                     string
	forChannelId                  string
	maxResults                    int64
	mine                          = utils.BoolPtr("false")
	myRecentSubscribers           = utils.BoolPtr("false")
	mySubscribers                 = utils.BoolPtr("false")
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	title                         string
	parts                         []string
	output                        string
	jpath                         string
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]**bool{
			"mine":                &mine,
			"myRecentSubscribers": &myRecentSubscribers,
			"mySubscribers":       &mySubscribers,
		}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(subscriptionCmd)
}

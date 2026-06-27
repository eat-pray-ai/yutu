// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short      = "Manage YouTube subscriptions"
	long       = "Manage YouTube subscriptions. Use this tool to list subscriptions/subscribers, subscribe to channels, or unsubscribe."
	scidUsage  = "Subscriber's channel id"
	descUsage  = "Description of the subscription"
	fcidUsage  = "Return the subscriptions to the subset of these channels that the authenticated user is subscribed to"
	forUsage = "mine|myRecentSubscribers|mySubscribers"
	orderUsage = "subscriptionOrderUnspecified|relevance|unread|alphabetical"
	titleUsage = "Title of the subscription"
)

var (
	ids                           []string
	subscriberChannelId           string
	description                   string
	channelId                     string
	forChannelId                  string
	maxResults                    int64
	subscriptionFor               string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	title                         string
	parts                         []string
)

var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(subscriptionCmd)
}

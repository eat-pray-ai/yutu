// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short         = "Manage YouTube channels"
	long          = "Manage YouTube channels. Use this tool to list or update channels."
	cidUsage      = "Return the channels within the specified guide category id"
	fhUsage       = "Return the channel associated with a YouTube handle"
	fuUsage       = "Return the channel associated with a YouTube username"
	hlUsage       = "Specifies the localization language of the metadata"
	forUsage      = "managedByMe|mine|mySubscribers"
	countryUsage  = "Country of the channel"
	keywordsUsage = "Comma separated channel keywords"
	dlUsage       = "The language of the channel's default title and description"
	descUsage     = "Description of the channel"
	titleUsage    = "Title of the channel"
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: short,
	Long:  long,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelCmd)
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short          = "Manage YouTube third-party links"
	long           = "Manage YouTube third-party links. Use this tool to list, insert, update, or delete links between a YouTube channel and a third-party service."
	ltUsage        = "Linking token that identifies the YouTube account and channel"
	typeUsage      = "Type of the link: linkUnspecified|channelToStoreLink"
	statusUsage    = "Status of the link: unknown|failed|pending|linked"
	extCidUsage    = "Channel ID to which changes should be applied, for delegation"
)

var (
	linkingToken      string
	linkType          string
	linkStatus        string
	externalChannelId string
	parts             []string
)

var thirdPartyLinkCmd = &cobra.Command{
	Use:   "thirdPartyLink",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(thirdPartyLinkCmd)
}
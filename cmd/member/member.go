// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "List channel's members' info"
	long      = "List channel's members' info, such as channelId, displayName, etc"
	mcidUsage = "Comma separated list of channel IDs. Only data about members that are part of this list will be included"
	hatlUsage = "Filter members in the results set to the ones that have access to a level"
	mmUsage   = "listMembersModeUnknown, updates, or all_current"
)

var (
	memberChannelId  string
	hasAccessToLevel string
	maxResults       int64
	mode             string
	parts            []string
	output           string
	jsonpath         string
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

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Manage YouTube live chat moderators"
	long      = "Manage YouTube live chat moderators. Use this tool to list, add, or remove moderators from a live chat."
	lcidUsage = "ID of the live chat"
	mcidUsage = "Channel ID of the moderator"
)

var (
	ids                []string
	liveChatId         string
	moderatorChannelId string
	maxResults         int64
	parts              []string
)

var liveChatModeratorCmd = &cobra.Command{
	Use:   "liveChatModerator",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveChatModeratorCmd)
}

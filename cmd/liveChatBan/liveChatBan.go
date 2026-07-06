// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short            = "Manage YouTube live chat bans"
	long             = "Manage YouTube live chat bans. Use this tool to ban or unban users from a live chat."
	lcidUsage        = "ID of the live chat"
	bucidUsage       = "Channel ID of the user to ban"
	banTypeUsage     = "permanent|temporary"
	banDurationUsage = "Duration of the ban in seconds (only for temporary bans)"
)

var (
	ids                 []string
	liveChatId          string
	bannedUserChannelId string
	banType             string
	banDurationSeconds  uint64
	parts               []string
)

var liveChatBanCmd = &cobra.Command{
	Use:   "liveChatBan",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveChatBanCmd)
}

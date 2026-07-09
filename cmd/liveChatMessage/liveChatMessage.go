// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "Manage YouTube live chat messages"
	long        = "Manage YouTube live chat messages. Use this tool to list, send, delete, or transition messages in a live chat."
	lcidUsage   = "ID of the live chat"
	msgUsage    = "Text of the message to send"
	statusUsage = "Status to transition to (closed)"
)

var (
	ids         []string
	liveChatId  string
	messageText string
	status      string
	maxResults  int64
	parts       []string
	hl          string
)

var liveChatMessageCmd = &cobra.Command{
	Use:   "liveChatMessage",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveChatMessageCmd)
}

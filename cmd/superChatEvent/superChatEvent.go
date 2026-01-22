// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short   = "List Super Chat events for a channel"
	long    = "List Super Chat events for a channel"
	hlUsage = "Return rendered funding amounts in specified language"
)

var (
	hl         string
	maxResults int64
	parts      []string
	output     string
	jsonpath   string
)

var superChatEventCmd = &cobra.Command{
	Use:   "superChatEvent",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(superChatEventCmd)
}

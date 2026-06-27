// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Manage activities on YouTube"
	long      = "Manage activities on YouTube. Use this tool to list channel activities."
	ciUsage   = "ID of the channel"
	forUsage = "home|mine"
	paUsage   = "Filter on activities published after this date"
	pbUsage   = "Filter on activities published before this date"
	rcUsage   = "Display the content as seen by viewers in this country"
)

var (
	channelId       string
	activityFor     string
	maxResults      int64
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

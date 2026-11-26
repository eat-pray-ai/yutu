// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short    = "Manage YouTube activities"
	long     = "Manage YouTube activities. Use this tool to list channel activities."
	ciUsage  = "ID of the channel"
	forUsage = "home|mine"
	paUsage  = "Filter on activities published after this date"
	pbUsage  = "Filter on activities published before this date"
	rcUsage  = "Display the content as seen by viewers in this country"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: short,
	Long:  long,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

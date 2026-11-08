// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short     = "Manage YouTube channel sections"
	long      = "Manage YouTube channel sections. Use this tool to list or delete channel sections."
	cidUsage  = "Return the ChannelSections owned by the specified channel id"
	hlUsage   = "Return content in specified language"
	mineUsage = "Return the ChannelSections owned by the authenticated user"
)

var (
	ids                    []string
	channelId              string
	hl                     string
	mine                   = new(false)
	onBehalfOfContentOwner string
	parts                  []string
)

var channelSectionCmd = &cobra.Command{
	Use:   "channelSection",
	Short: short,
	Long:  long,
	PersistentPreRun: func(c *cobra.Command, _ []string) {
		utils.ResetFlags(map[string]**bool{"mine": &mine}, c.Flags())
	},
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelSectionCmd)
}

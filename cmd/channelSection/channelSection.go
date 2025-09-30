// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short     = "Manipulate YouTube channel sections"
	long      = "List or delete YouTube channel sections"
	cidUsage  = "Return the ChannelSections owned by the specified channel id"
	hlUsage   = "Return content in specified language"
	mineUsage = "Return the ChannelSections owned by the authenticated user"
)

var (
	ids                    []string
	channelId              string
	hl                     string
	mine                   = utils.BoolPtr("false")
	onBehalfOfContentOwner string
	parts                  []string
	output                 string
	jpath                  string
)

var channelSectionCmd = &cobra.Command{
	Use:   "channelSection",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ResetBool(map[string]**bool{"mine": &mine}, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelSectionCmd)
}

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
	long      = "List or delete YouTube channel sections\n\nExamples:\n  yutu channelSection list --mine\n  yutu channelSection list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw\n  yutu channelSection delete --ids abc123\n  yutu channelSection delete --ids abc123,def456"
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
	output                 string
	jsonpath               string
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Insert YouTube channel banner"
	long      = "Insert YouTube channel banner\n\nExamples:\n  yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.jpg\n  yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.png --output json"
	cidUsage  = "ID of the channel to insert the banner for"
	fileUsage = "Path to the banner image"
)

var (
	channelId string
	file      string
	output    string
	jsonpath  string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var channelBannerCmd = &cobra.Command{
	Use:   "channelBanner",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelBannerCmd)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Manage YouTube watermarks"
	long      = "Manage YouTube watermarks. Use this tool to set or unset watermarks for channel videos."
	cidUsage  = "ID of the channel"
	fileUsage = "Path to the watermark file"
	ivpUsage  = "topLeft|topRight|bottomLeft|bottomRight"
	dmUsage   = "Duration in milliseconds for which the watermark should be displayed"
	omUsage   = "Defines the time at which the watermark will appear"
	otUsage   = "offsetFromStart|offsetFromEnd"
)

var (
	channelId              string
	file                   string
	inVideoPosition        string
	durationMs             uint64
	offsetMs               uint64
	offsetType             string
	onBehalfOfContentOwner string
)

var watermarkCmd = &cobra.Command{
	Use:   "watermark",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(watermarkCmd)
}

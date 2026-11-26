// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "Manage YouTube live streams"
	long        = "Manage YouTube live streams. Use this tool to list, insert, update, or delete live streams."
	titleUsage  = "Title of the live stream"
	descUsage   = "Description of the live stream"
	frUsage     = "Frame rate of the video stream (30fps, 60fps, variable)"
	itUsage     = "Ingestion type (rtmp, dash, webrtc, hls)"
	resUsage    = "Resolution of the video stream (240p, 360p, 480p, 720p, 1080p, 1440p, 2160p, variable)"
	obococUsage = "Channel ID for content owner operations"
)

var liveStreamCmd = &cobra.Command{
	Use:   "liveStream",
	Short: short,
	Long:  long,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveStreamCmd)
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

var (
	ids                           []string
	title                         string
	description                   string
	mine                          = new(false)
	frameRate                     string
	ingestionType                 string
	resolution                    string
	maxResults                    int64
	parts                         []string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var liveStreamCmd = &cobra.Command{
	Use:   "liveStream",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resetFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveStreamCmd)
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]**bool{
		"mine": &mine,
	}

	utils.ResetBool(boolMap, flagSet)
}

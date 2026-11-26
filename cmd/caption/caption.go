// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short         = "Manage YouTube video captions"
	long          = "Manage YouTube video captions. Use this tool to list, insert, update, download, or delete video captions."
	fileUsage     = "Path to save the caption file"
	attUsage      = "unknown|primary|commentary|descriptive"
	iasUsage      = "Whether YouTube synchronized the caption track to the audio track in the video"
	iscUsage      = "Whether the track contains closed captions for the deaf and hard of hearing"
	isdUsage      = "whether the caption track is a draft"
	iserUsage     = "Whether caption track is formatted for 'easy reader'"
	islUsage      = "Whether the caption track uses large text for the vision-impaired"
	languageUsage = "Language of the caption track"
	nameUsage     = "Name of the caption track"
	tkUsage       = "standard|ASR|forced"
	vidUsage      = "ID of the video"
	tfmtUsage     = "sbv|srt|vtt"
	tlangUsage    = "Translate the captions into this language"
)

var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: short,
	Long:  long,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(captionCmd)
}

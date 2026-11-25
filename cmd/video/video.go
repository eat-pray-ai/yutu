// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

const (
	short           = "Manage YouTube videos"
	long            = "Manage YouTube videos. Use this tool to list, upload, update, delete, get rating, or report videos."
	alUsage         = "Should auto-levels be applied to the upload"
	fileUsage       = "Path to the video file"
	titleUsage      = "Title of the video"
	descUsage       = "Description of the video"
	hlUsage         = "Specifies the localization language"
	tagsUsage       = "Comma separated tags"
	localeUsage     = "Locale of the video metadata to return"
	licenseUsage    = "youtube|creativeCommon"
	thumbnailUsage  = "Path to the thumbnail file"
	chartUsage      = "chartUnspecified|mostPopular"
	chidUsage       = "Channel id of the video"
	commentsUsage   = "Additional comments regarding the abuse report"
	pidUsage        = "Playlist id of the video"
	caidUsage       = "Category of the video"
	privacyUsage    = "public|private|unlisted"
	fkUsage         = "Whether the video is for kids"
	embeddableUsage = "Whether the video is embeddable"
	csmUsage        = "Whether the video contains altered or synthetic media"
	rdUsage         = "Date and time when the video was recorded"
	paUsage         = "Datetime when the video is scheduled to publish"
	rcUsage         = "Specific to the specified region"
	ridUsage        = "ID of the reason for reporting abuse"
	sridUsage       = "ID of the secondary reason for reporting abuse"
	stabilizeUsage  = "Should stabilize be applied to the upload"
	mhUsage         = "Max height of the embedded player in pixels"
	mwUsage         = "Max width of the embedded player in pixels"
	nsUsage         = "Notify the channel subscribers about the new video"
	psvUsage        = "Whether the extended video statistics can be viewed by everyone"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: short,
	Long:  long,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCmd)
}

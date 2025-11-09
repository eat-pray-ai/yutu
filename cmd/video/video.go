// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

const (
	short           = "Manipulate YouTube videos"
	long            = "List, insert, update, rate, get rating, report abuse, or delete YouTube videos"
	alUsage         = "Should auto-levels be applied to the upload"
	fileUsage       = "Path to the video file"
	titleUsage      = "Title of the video"
	descUsage       = "Description of the video"
	hlUsage         = "Specifies the localization language"
	tagsUsage       = "Comma separated tags"
	localUsage      = ""
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
	paUsage         = "Datetime when the video is scheduled to publish"
	rcUsage         = "Specific to the specified region"
	ridUsage        = "ID of the reason for reporting abuse"
	sridUsage       = "ID of the secondary reason for reporting abuse"
	stabilizeUsage  = "Should stabilize be applied to the upload"
	mhUsage         = ""
	mwUsage         = ""
	nsUsage         = "Notify the channel subscribers about the new video"
	psvUsage        = "Whether the extended video statistics can be viewed by everyone"
)

var (
	ids               []string
	autoLevels        = utils.BoolPtr("false")
	file              string
	title             string
	description       string
	hl                string
	tags              []string
	language          string
	locale            string
	license           string
	thumbnail         string
	rating            string
	chart             string
	channelId         string
	comments          string
	playListId        string
	categoryId        string
	privacy           string
	forKids           = utils.BoolPtr("false")
	embeddable        = utils.BoolPtr("false")
	publishAt         string
	regionCode        string
	reasonId          string
	secondaryReasonId string
	stabilize         = utils.BoolPtr("false")
	maxHeight         int64
	maxWidth          int64
	maxResults        int64
	parts             []string
	output            string
	jpath             string

	notifySubscribers             = utils.BoolPtr("false")
	publicStatsViewable           = utils.BoolPtr("false")
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var videoCmd = &cobra.Command{
	Use:   "video",
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
	cmd.RootCmd.AddCommand(videoCmd)
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]**bool{
		"autoLevels":          &autoLevels,
		"forKids":             &forKids,
		"embeddable":          &embeddable,
		"stabilize":           &stabilize,
		"notifySubscribers":   &notifySubscribers,
		"publicStatsViewable": &publicStatsViewable,
	}

	utils.ResetBool(boolMap, flagSet)
}

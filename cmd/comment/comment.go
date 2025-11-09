// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short     = "Manipulate YouTube comments"
	long      = "List, insert, update, mark as spam, set moderation status, or delete YouTube comments"
	idsUsage  = "IDs of comments"
	acidUsage = "Channel id of the comment author"
	crUsage   = "Whether the viewer can rate the comment"
	cidUsage  = "Channel id of the video owner"
	tfUsage   = "textFormatUnspecified|html|plainText"
	toUsage   = "Text of the comment"
	msUsage   = "heldForReview|published|rejected"
	baUsage   = "If set to true the author of the comment gets added to the ban list"
	vidUsage  = "ID of the video"
	vrUsage   = "none|like|dislike"
)

var (
	ids              []string
	authorChannelId  string
	canRate          = utils.BoolPtr("false")
	channelId        string
	maxResults       int64
	parentId         string
	textFormat       string
	textOriginal     string
	moderationStatus string
	banAuthor        = utils.BoolPtr("false")
	videoId          string
	viewerRating     string
	parts            []string
	output           string
	jpath            string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]**bool{"canRate": &canRate, "banAuthor": &banAuthor}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentCmd)
}

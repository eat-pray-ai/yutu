// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "comment-insert"
	insertShort    = "Insert a comment"
	insertLong     = "Insert a comment to a video\n\nExamples:\n  yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Hello'\n  yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Reply' --canRate"
	insertPidUsage = "ID of the parent comment"
)

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"author_channel_id", "channel_id", "parent_id", "text_original", "video_id",
	},
	Properties: map[string]*jsonschema.Schema{
		"author_channel_id": {Type: "string", Description: acidUsage},
		"channel_id":        {Type: "string", Description: cidUsage},
		"can_rate":          {Type: "boolean", Description: crUsage},
		"parent_id":         {Type: "string", Description: insertPidUsage},
		"text_original":     {Type: "string", Description: toUsage},
		"video_id":          {Type: "string", Description: vidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			insertTool, func(input comment.Comment, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	commentCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", cidUsage,
	)
	insertCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	insertCmd.Flags().StringVarP(
		&parentId, "parentId", "P", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("parentId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := comment.NewComment(
			comment.WithAuthorChannelId(authorChannelId),
			comment.WithChannelId(channelId),
			comment.WithCanRate(canRate),
			comment.WithParentId(parentId),
			comment.WithTextOriginal(textOriginal),
			comment.WithVideoId(videoId),
			comment.WithOutput(output),
			comment.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Insert(cmd.OutOrStdout()), cmd)
	},
}

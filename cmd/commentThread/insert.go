// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "commentThread-insert"
	insertShort    = "Insert a new comment thread"
	insertLong     = "Insert a new comment thread"
	insertVidUsage = "ID of the video"
)

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"author_channel_id", "channel_id", "text_original", "video_id",
	},
	Properties: map[string]*jsonschema.Schema{
		"author_channel_id": {Type: "string", Description: acidUsage},
		"channel_id":        {Type: "string", Description: cidUsage},
		"text_original":     {Type: "string", Description: toUsage},
		"video_id":          {Type: "string", Description: insertVidUsage},
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
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			insertTool,
			func(input commentThread.CommentThread, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", toUsage)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", insertVidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := commentThread.NewCommentThread(
			commentThread.WithAuthorChannelId(authorChannelId),
			commentThread.WithChannelId(channelId),
			commentThread.WithTextOriginal(textOriginal),
			commentThread.WithVideoId(videoId),
			commentThread.WithOutput(output),
			commentThread.WithJsonpath(jsonpath),
		)
		err := input.Insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

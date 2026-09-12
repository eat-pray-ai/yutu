// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "commentThread-insert"
	insertVidUsage = "ID of the video"
	insertConfirm  = "Create comment thread on video: %s"
	insertShort    = "Create a comment thread"
	insertLong     = "Create a comment thread. Use this tool to create a comment thread."
	insertExample  = `# Post a comment on a video
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Great video!'
# Post a comment with JSON output
yutu commentThread insert --channelId UC_x5X --videoId dQw4w9WgXcQ --authorChannelId UA_x5X --textOriginal 'Nice work!' --output json`
)

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"channel_id", "text_original", "video_id",
	},
	Properties: map[string]*jsonschema.Schema{
		"author_channel_id": {Type: "string", Description: acidUsage},
		"channel_id":        {Type: "string", Description: cidUsage},
		"text_original":     {Type: "string", Description: toUsage},
		"video_id":          {Type: "string", Description: insertVidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
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
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input commentThread.CommentThread) string {
					return fmt.Sprintf(insertConfirm, input.VideoId)
				},
				func(input commentThread.CommentThread, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", toUsage)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", insertVidUsage)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(insertConfirm, videoId))
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := commentThread.NewCommentThread(
			commentThread.WithAuthorChannelId(authorChannelId),
			commentThread.WithChannelId(channelId),
			commentThread.WithTextOriginal(textOriginal),
			commentThread.WithVideoId(videoId),
			commentThread.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

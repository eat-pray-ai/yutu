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
	insertConfirm  = "Create comment thread on video %s"
	insertShort    = "Create a comment thread"
	insertLong     = "Create a comment thread. Use this tool to create a comment thread."
	insertExample  = `  # Post a comment on a video
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
	commentThreadCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	authorChannelId string
	channelId       string
	textOriginal    string
	videoId         string
	output          string
	yes             bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&o.authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().StringVarP(&o.textOriginal, "textOriginal", "t", "", toUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", insertVidUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("channelId")
	_ = c.MarkFlagRequired("textOriginal")
	_ = c.MarkFlagRequired("videoId")
}

func (o *insertOpts) build() *commentThread.CommentThread {
	return commentThread.NewCommentThread(
		commentThread.WithAuthorChannelId(o.authorChannelId),
		commentThread.WithChannelId(o.channelId),
		commentThread.WithTextOriginal(o.textOriginal),
		commentThread.WithVideoId(o.videoId),
		commentThread.WithOutput(o.output),
	).(*commentThread.CommentThread)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.videoId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

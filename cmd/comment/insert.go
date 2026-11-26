// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	insertTool     = "comment-insert"
	insertPidUsage = "ID of the parent comment"
	insertConfirm  = "Create comment on parent %s"
	insertShort    = "Create a comment"
	insertLong     = "Create a comment. Use this tool to create a comment on a video."
	insertExample  = `  # Reply to a comment
  yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Hello'
  # Reply with rating enabled
  yutu comment insert --channelId UC_x5X --videoId dQw4w9 --authorChannelId UA_x5X --parentId UgyXXX --textOriginal 'Reply' --canRate`
)

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"parent_id", "text_original",
	},
	Properties: map[string]*jsonschema.Schema{
		"author_channel_id": {Type: "string", Description: acidUsage},
		"channel_id":        {Type: "string", Description: cidUsage},
		"can_rate":          {Type: "boolean", Description: crUsage},
		"parent_id":         {Type: "string", Description: insertPidUsage},
		"text_original":     {Type: "string", Description: toUsage},
		"video_id":          {Type: "string", Description: vidUsage},
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
				func(input comment.Comment) string {
					return fmt.Sprintf(insertConfirm, input.ParentId)
				},
				func(input comment.Comment, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	commentCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	authorChannelId string
	channelId       string
	canRate         bool
	parentId        string
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
	c.Flags().BoolVarP(&o.canRate, "canRate", "R", false, crUsage)
	c.Flags().StringVarP(&o.parentId, "parentId", "P", "", insertPidUsage)
	c.Flags().StringVarP(&o.textOriginal, "textOriginal", "t", "", toUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("parentId")
	_ = c.MarkFlagRequired("textOriginal")
}

func (o *insertOpts) build(f *pflag.FlagSet) *comment.Comment {
	return comment.NewComment(
		comment.WithAuthorChannelId(o.authorChannelId),
		comment.WithChannelId(o.channelId),
		utils.OptIfChanged(f, "canRate", &o.canRate, comment.WithCanRate),
		comment.WithParentId(o.parentId),
		comment.WithTextOriginal(o.textOriginal),
		comment.WithVideoId(o.videoId),
		comment.WithOutput(o.output),
	).(*comment.Comment)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.parentId))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

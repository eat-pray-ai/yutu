// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatMessage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "liveChatMessage-insert"
	insertConfirm = "Send message to live chat %s"
	insertShort   = "Send a live chat message"
	insertLong    = "Send a live chat message. Use this tool to post a text message to a live chat."
	insertExample = `  # Send a message to a live chat
  yutu liveChatMessage insert --liveChatId abc123 --messageText "Hello everyone!"`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id", "message_text"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id": {Type: "string", Description: lcidUsage},
		"message_text": {Type: "string", Description: msgUsage},
		"parts": {
			Type: "array", Description: "Parts to include in the response",
			Items: &jsonschema.Schema{Type: "string"}, Default: jsontext.Value(`["snippet"]`),
		},
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
				func(input liveChatMessage.LiveChatMessage) string {
					return fmt.Sprintf(insertConfirm, input.LiveChatId)
				},
				func(input liveChatMessage.LiveChatMessage, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	liveChatMessageCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	liveChatId  string
	messageText string
	parts       []string
	output      string
	yes         bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.liveChatId, "liveChatId", "L", "", lcidUsage)
	c.Flags().StringVarP(&o.messageText, "messageText", "m", "", msgUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("liveChatId")
	_ = c.MarkFlagRequired("messageText")
}

func (o *insertOpts) build() *liveChatMessage.LiveChatMessage {
	return liveChatMessage.NewLiveChatMessage(
		liveChatMessage.WithLiveChatId(o.liveChatId),
		liveChatMessage.WithMessageText(o.messageText),
		liveChatMessage.WithParts(o.parts),
		liveChatMessage.WithOutput(o.output),
	).(*liveChatMessage.LiveChatMessage)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.liveChatId))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

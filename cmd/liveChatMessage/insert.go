// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"encoding/json"
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
	insertShort   = "Send a live chat message"
	insertLong    = "Send a live chat message. Use this tool to post a text message to a live chat."
	insertExample = `# Send a message to a live chat
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
			Items: &jsonschema.Schema{Type: "string"}, Default: json.RawMessage(`["snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
		}, cobramcp.GenToolHandler(
			insertTool,
			func(input liveChatMessage.LiveChatMessage, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	liveChatMessageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&liveChatId, "liveChatId", "l", "", lcidUsage)
	insertCmd.Flags().StringVarP(&messageText, "messageText", "m", "", msgUsage)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("liveChatId")
	_ = insertCmd.MarkFlagRequired("messageText")
	cmd.AddMutationFlags(insertCmd)
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would send message to live chat %s", liveChatId,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveChatMessage.NewLiveChatMessage(
			liveChatMessage.WithLiveChatId(liveChatId),
			liveChatMessage.WithMessageText(messageText),
			liveChatMessage.WithParts(parts),
			liveChatMessage.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

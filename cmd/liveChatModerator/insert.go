// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"encoding/json"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatModerator"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "liveChatModerator-insert"
	insertShort   = "Insert a live chat moderator"
	insertLong    = "Insert a live chat moderator. Use this tool to add a moderator to a live chat."
	insertExample = `# Add a moderator to a live chat
yutu liveChatModerator insert --liveChatId abc123 --moderatorChannelId UC_xyz`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id", "moderator_channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id":         {Type: "string", Description: lcidUsage},
		"moderator_channel_id": {Type: "string", Description: mcidUsage},
		"parts": {
			Type: "array", Description: "Parts to include in the response",
			Items: &jsonschema.Schema{Type: "string"}, Default: json.RawMessage(`["snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
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
			func(input liveChatModerator.LiveChatModerator, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Insert(writer)
			},
		),
	)
	liveChatModeratorCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&liveChatId, "liveChatId", "l", "", lcidUsage)
	insertCmd.Flags().StringVarP(
		&moderatorChannelId, "moderatorChannelId", "m", "", mcidUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("liveChatId")
	_ = insertCmd.MarkFlagRequired("moderatorChannelId")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf(
			"Would add moderator %s to live chat %s", moderatorChannelId, liveChatId,
		)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := liveChatModerator.NewLiveChatModerator(
			liveChatModerator.WithLiveChatId(liveChatId),
			liveChatModerator.WithModeratorChannelId(moderatorChannelId),
			liveChatModerator.WithParts(parts),
			liveChatModerator.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

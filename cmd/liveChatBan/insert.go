// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"encoding/json"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatBan"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "liveChatBan-insert"
	insertShort   = "Insert a live chat ban"
	insertLong    = "Insert a live chat ban. Use this tool to ban a user from a live chat."
	insertExample = `# Ban a user permanently
yutu liveChatBan insert --liveChatId abc123 --bannedUserChannelId UC_xyz --banType permanent
# Ban a user temporarily for 5 minutes
yutu liveChatBan insert --liveChatId abc123 --bannedUserChannelId UC_xyz --banType temporary --banDurationSeconds 300`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id", "banned_user_channel_id", "ban_type"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id":           {Type: "string", Description: lcidUsage},
		"banned_user_channel_id": {Type: "string", Description: bucidUsage},
		"ban_type":               {Type: "string", Description: banTypeUsage},
		"ban_duration_seconds":   {Type: "number", Description: banDurationUsage},
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
			insertTool, func(input liveChatBan.LiveChatBan, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	liveChatBanCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&liveChatId, "liveChatId", "l", "", lcidUsage)
	insertCmd.Flags().StringVarP(
		&bannedUserChannelId, "bannedUserChannelId", "b", "", bucidUsage,
	)
	insertCmd.Flags().StringVarP(&banType, "banType", "t", "", banTypeUsage)
	insertCmd.Flags().Uint64VarP(
		&banDurationSeconds, "banDurationSeconds", "d", 0, banDurationUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, "Parts to include",
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("liveChatId")
	_ = insertCmd.MarkFlagRequired("bannedUserChannelId")
	_ = insertCmd.MarkFlagRequired("banType")
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
			c, "Would ban user %s in live chat %s", bannedUserChannelId, liveChatId,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveChatBan.NewLiveChatBan(
			liveChatBan.WithLiveChatId(liveChatId),
			liveChatBan.WithBannedUserChannelId(bannedUserChannelId),
			liveChatBan.WithBanType(banType),
			liveChatBan.WithBanDurationSeconds(banDurationSeconds),
			liveChatBan.WithParts(parts),
			liveChatBan.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

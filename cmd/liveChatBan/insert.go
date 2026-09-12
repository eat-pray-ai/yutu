// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"encoding/json/jsontext"
	"fmt"
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
	insertConfirm = "Ban user %s in live chat %s"
	insertShort   = "Create a live chat ban"
	insertLong    = "Create a live chat ban. Use this tool to ban a user from a live chat."
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
		"ban_type": {
			Type: "string", Description: banTypeUsage,
			Enum: []any{"permanent", "temporary"},
		},
		"ban_duration_seconds": {
			Type: "number", Description: banDurationUsage, Minimum: new(float64(0)),
		},
		"parts": {
			Type: "array", Description: "Parts to include in the response",
			Items: &jsonschema.Schema{Type: "string"}, Default: jsontext.Value(`["snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"yaml"`),
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
				func(input liveChatBan.LiveChatBan) string {
					return fmt.Sprintf(
						insertConfirm, input.BannedUserChannelId, input.LiveChatId,
					)
				},
				func(input liveChatBan.LiveChatBan, w io.Writer) error {
					return input.Insert(w)
				},
			),
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
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("liveChatId")
	_ = insertCmd.MarkFlagRequired("bannedUserChannelId")
	_ = insertCmd.MarkFlagRequired("banType")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(
			c, fmt.Sprintf(insertConfirm, bannedUserChannelId, liveChatId),
		)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
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

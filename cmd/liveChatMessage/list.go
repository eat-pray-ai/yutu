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
	listTool    = "liveChatMessage-list"
	listShort   = "List live chat messages"
	listLong    = "List live chat messages. Use this tool to list messages in a live chat."
	listExample = `# List messages for a live chat
yutu liveChatMessage list --liveChatId abc123
# List messages with limit
yutu liveChatMessage list --liveChatId abc123 --maxResults 10`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id": {Type: "string", Description: lcidUsage},
		"hl":           {Type: "string", Description: "Localization language for system messages"},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: new(float64(0)),
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["snippet","authorDetails"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cobramcp.GenToolHandler(
			listTool,
			func(input liveChatMessage.LiveChatMessage, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveChatMessageCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&liveChatId, "liveChatId", "l", "", lcidUsage)
	listCmd.Flags().StringVar(
		&hl, "hl", "", "Localization language for system messages",
	)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "authorDetails"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringP("output", "o", "table", pkg.TableUsage)

	_ = listCmd.MarkFlagRequired("liveChatId")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Long:    listLong,
	Example: listExample,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		input := liveChatMessage.NewLiveChatMessage(
			liveChatMessage.WithLiveChatId(liveChatId),
			liveChatMessage.WithHl(hl),
			liveChatMessage.WithMaxResults(maxResults),
			liveChatMessage.WithParts(parts),
			liveChatMessage.WithOutput(output),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

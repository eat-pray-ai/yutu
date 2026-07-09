// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/liveChatMessage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "liveChatMessage-delete"
	deleteIdsUsage = "IDs of the live chat messages to delete"
	deleteShort    = "Delete live chat messages"
	deleteLong     = "Delete live chat messages. Use this tool to remove messages from a live chat by their IDs."
	deleteExample  = `# Delete a live chat message by ID
yutu liveChatMessage delete --ids abc123
# Delete multiple live chat messages
yutu liveChatMessage delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(true),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			deleteTool,
			func(input liveChatMessage.LiveChatMessage, writer io.Writer) error {
				return input.Delete(writer)
			},
		),
	)
	liveChatMessageCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
	cmd.AddMutationFlags(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Long:    deleteLong,
	Example: deleteExample,
	Run: func(c *cobra.Command, args []string) {
		err := cmd.Confirm(
			c, "Would delete live chat message(s): %s", strings.Join(ids, ", "),
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveChatMessage.NewLiveChatMessage(
			liveChatMessage.WithIds(ids),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

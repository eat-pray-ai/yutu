// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/liveChatModerator"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "liveChatModerator-delete"
	deleteIdsUsage = "IDs of the live chat moderators to delete"
	deleteShort    = "Delete live chat moderators"
	deleteLong     = "Delete live chat moderators. Use this tool to remove moderators from a live chat by their IDs."
	deleteExample  = `# Delete a live chat moderator by ID
yutu liveChatModerator delete --ids abc123
# Delete multiple live chat moderators
yutu liveChatModerator delete --ids abc123,def456`
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
			func(input liveChatModerator.LiveChatModerator, writer io.Writer) error {
				return input.Delete(writer)
			},
		),
	)
	liveChatModeratorCmd.AddCommand(deleteCmd)

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
			c, "Would delete live chat moderator(s): %s", strings.Join(ids, ", "),
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveChatModerator.NewLiveChatModerator(
			liveChatModerator.WithIds(ids),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

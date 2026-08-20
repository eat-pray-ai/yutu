// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"fmt"
	"io"
	"strings"

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
	deleteTool     = "liveChatBan-delete"
	deleteIdsUsage = "IDs of the live chat bans to delete"
	deleteShort    = "Delete live chat bans"
	deleteLong     = "Delete live chat bans. Use this tool to unban users from a live chat by ban IDs."
	deleteExample  = `# Delete a live chat ban by ID
yutu liveChatBan delete --ids abc123
# Delete multiple live chat bans
yutu liveChatBan delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
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
			deleteTool, func(input liveChatBan.LiveChatBan, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Delete(writer)
			},
		),
	)
	liveChatBanCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Long:    deleteLong,
	Example: deleteExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf(
			"Would delete live chat ban(s): %s", strings.Join(ids, ", "),
		)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		input := liveChatBan.NewLiveChatBan(
			liveChatBan.WithIds(ids),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

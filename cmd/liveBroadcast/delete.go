// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveBroadcast"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "liveBroadcast-delete"
	deleteIdsUsage = "IDs of the live broadcasts to delete"
	deleteShort    = "Delete live broadcasts"
	deleteLong     = "Delete live broadcasts. Use this tool to delete live broadcasts by their IDs."
	deleteExample  = `# Delete a live broadcast by ID
yutu liveBroadcast delete --ids broadcast123
# Delete multiple live broadcasts
yutu liveBroadcast delete --ids broadcast123,broadcast456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
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
			func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.Delete(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
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
			c, "Would delete live broadcast(s): %s", strings.Join(ids, ", "),
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

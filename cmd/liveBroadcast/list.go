// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json"
	"io"

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
	listTool     = "liveBroadcast-list"
	listIdsUsage = "Return live broadcasts with the given ids"
	listShort    = "List live broadcasts"
	listLong     = "List live broadcasts. Use this tool to list live broadcasts for the authenticated user."
	listExample  = `# List my live broadcasts
yutu liveBroadcast list --mine
# List live broadcasts by ID
yutu liveBroadcast list --ids broadcast1,broadcast2
# List active live broadcasts
yutu liveBroadcast list --mine --broadcastStatus active --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"mine":             {Type: "boolean", Description: "List broadcasts owned by the authenticated user"},
		"broadcast_status": {Type: "string", Enum: []any{"all", "active", "upcoming", "completed"}, Description: bsUsage},
		"broadcast_type":   {Type: "string", Enum: []any{"all", "event", "persistent"}, Description: btUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"), Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet","status"]`),
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
			listTool, func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().BoolVarP(
		mine, "mine", "m", false, "List broadcasts owned by the authenticated user",
	)
	listCmd.Flags().StringVarP(
		&broadcastStatus, "broadcastStatus", "s", "", bsUsage,
	)
	listCmd.Flags().StringVarP(&broadcastType, "broadcastType", "T", "", btUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"},
		pkg.PartsUsage,
	)
	listCmd.Flags().StringP("output", "o", "table", pkg.TableUsage)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Long:    listLong,
	Example: listExample,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithMine(mine),
			liveBroadcast.WithBroadcastStatus(broadcastStatus),
			liveBroadcast.WithBroadcastType(broadcastType),
			liveBroadcast.WithMaxResults(maxResults),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithParts(parts),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

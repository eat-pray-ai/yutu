// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json"
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
	transitionTool    = "liveBroadcast-transition"
	transitionBsUsage = "Broadcast status to transition to (testing, live, complete)"
	transitionShort   = "Transition a live broadcast"
	transitionLong    = "Transition a live broadcast. Use this tool to change the status of a live broadcast (e.g., go live, end broadcast)."
	transitionExample = `# Transition a broadcast to live
yutu liveBroadcast transition --ids broadcast123 --broadcastStatus live
# End a broadcast
yutu liveBroadcast transition --ids broadcast123 --broadcastStatus complete`
)

var transitionInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "broadcast_status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live broadcasts to transition",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"broadcast_status": {
			Type: "string", Enum: []any{"testing", "live", "complete"},
			Description: transitionBsUsage,
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
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: transitionTool, Title: transitionShort, Description: transitionLong,
			InputSchema: transitionInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			transitionTool,
			func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.Transition(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(transitionCmd)

	transitionCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, "IDs of the broadcasts to transition",
	)
	transitionCmd.Flags().StringVarP(
		&broadcastStatus, "broadcastStatus", "s", "", transitionBsUsage,
	)
	transitionCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	transitionCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	transitionCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	transitionCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = transitionCmd.MarkFlagRequired("ids")
	_ = transitionCmd.MarkFlagRequired("broadcastStatus")
	cmd.AddMutationFlags(transitionCmd)
}

var transitionCmd = &cobra.Command{
	Use:     "transition",
	Short:   transitionShort,
	Long:    transitionLong,
	Example: transitionExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would transition live broadcast(s) %s to %s",
			strings.Join(ids, ", "), broadcastStatus,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithBroadcastStatus(broadcastStatus),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithParts(parts),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.Transition(c.OutOrStdout()), c)
	},
}

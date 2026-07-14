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
	bindTool    = "liveBroadcast-bind"
	bindShort   = "Bind a live broadcast to a stream"
	bindLong    = "Bind a live broadcast to a stream. Use this tool to bind or unbind a live stream to/from a live broadcast."
	bindExample = `# Bind a broadcast to a stream
yutu liveBroadcast bind --ids broadcast123 --streamId stream456
# Unbind a broadcast (omit streamId)
yutu liveBroadcast bind --ids broadcast123`
)

var bindInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live broadcasts to bind",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"stream_id":                  {Type: "string", Description: sidUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet","contentDetails","status"]`),
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
			Name: bindTool, Title: bindShort, Description: bindLong,
			InputSchema: bindInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			bindTool, func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.Bind(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(bindCmd)

	bindCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, "IDs of the broadcasts to bind",
	)
	bindCmd.Flags().StringVarP(&streamId, "streamId", "s", "", sidUsage)
	bindCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	bindCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	bindCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "contentDetails", "status"},
		pkg.PartsUsage,
	)
	bindCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = bindCmd.MarkFlagRequired("ids")
	cmd.AddMutationFlags(bindCmd)
}

var bindCmd = &cobra.Command{
	Use:     "bind",
	Short:   bindShort,
	Long:    bindLong,
	Example: bindExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would bind live broadcast(s) %s to stream %s",
			strings.Join(ids, ", "), streamId,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithStreamId(streamId),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithParts(parts),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.Bind(c.OutOrStdout()), c)
	},
}

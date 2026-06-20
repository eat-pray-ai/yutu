// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"encoding/json"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thirdPartyLink"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "thirdPartyLink-list"
	listShort   = "List third-party links"
	listLong    = "List third-party links. Use this tool to list links between a YouTube channel and a third-party service."
	listExample = `# List all third-party links
yutu thirdPartyLink list --parts snippet,status
# List links by linking token
yutu thirdPartyLink list --linkingToken abc123 --parts snippet,status
# List links by type
yutu thirdPartyLink list --type channelToStoreLink --parts snippet,status`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"linking_token":      {Type: "string", Description: ltUsage},
		"type":               {Type: "string", Description: typeUsage, Enum: []any{"linkUnspecified", "channelToStoreLink"}},
		"external_channel_id": {Type: "string", Description: extCidUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["snippet","status"]`),
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
			listTool, func(input thirdPartyLink.ThirdPartyLink, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	thirdPartyLinkCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&linkingToken, "linkingToken", "l", "", ltUsage)
	listCmd.Flags().StringVarP(&linkType, "type", "t", "", typeUsage)
	listCmd.Flags().StringVarP(&externalChannelId, "externalChannelId", "e", "", extCidUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
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
		input := thirdPartyLink.NewThirdPartyLink(
			thirdPartyLink.WithLinkingToken(linkingToken),
			thirdPartyLink.WithType(linkType),
			thirdPartyLink.WithExternalChannelId(externalChannelId),
			thirdPartyLink.WithParts(parts),
			thirdPartyLink.WithOutput(output),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}
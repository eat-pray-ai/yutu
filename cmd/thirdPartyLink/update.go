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
	updateTool    = "thirdPartyLink-update"
	updateShort   = "Update a third-party link"
	updateLong    = "Update a third-party link. Use this tool to update the status or type of a link between a YouTube channel and a third-party service."
	updateExample = `# Update a third-party link status
yutu thirdPartyLink update --linkingToken abc123 --linkStatus linked --parts snippet,status
# Update a third-party link type
yutu thirdPartyLink update --linkingToken abc123 --type channelToStoreLink --parts snippet,status`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"linking_token"},
	Properties: map[string]*jsonschema.Schema{
		"linking_token":       {Type: "string", Description: ltUsage},
		"type":                {Type: "string", Description: typeUsage, Enum: []any{"linkUnspecified", "channelToStoreLink"}},
		"link_status":         {Type: "string", Description: statusUsage, Enum: []any{"unknown", "failed", "pending", "linked"}},
		"external_channel_id": {Type: "string", Description: extCidUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["snippet","status"]`),
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			updateTool, func(input thirdPartyLink.ThirdPartyLink, writer io.Writer) error {
				return input.Update(writer)
			},
		),
	)
	thirdPartyLinkCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&linkingToken, "linkingToken", "l", "", ltUsage)
	updateCmd.Flags().StringVarP(&linkType, "type", "t", "", typeUsage)
	updateCmd.Flags().StringVarP(&linkStatus, "linkStatus", "s", "", statusUsage)
	updateCmd.Flags().StringVarP(&externalChannelId, "externalChannelId", "e", "", extCidUsage)
	updateCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
	)
	updateCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = updateCmd.MarkFlagRequired("linkingToken")
	cmd.AddMutationFlags(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(c, "Would update third-party link: %s", linkingToken)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := thirdPartyLink.NewThirdPartyLink(
			thirdPartyLink.WithLinkingToken(linkingToken),
			thirdPartyLink.WithType(linkType),
			thirdPartyLink.WithLinkStatus(linkStatus),
			thirdPartyLink.WithExternalChannelId(externalChannelId),
			thirdPartyLink.WithParts(parts),
			thirdPartyLink.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}
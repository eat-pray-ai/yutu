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
	insertTool    = "thirdPartyLink-insert"
	insertShort   = "Insert a new third-party link"
	insertLong    = "Insert a new third-party link. Use this tool to create a link between a YouTube channel and a third-party service."
	insertExample = `# Insert a new third-party link
yutu thirdPartyLink insert --linkingToken abc123 --type channelToStoreLink --linkStatus pending --parts snippet,status`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"linking_token", "type"},
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			insertTool, func(input thirdPartyLink.ThirdPartyLink, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	thirdPartyLinkCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&linkingToken, "linkingToken", "l", "", ltUsage)
	insertCmd.Flags().StringVarP(&linkType, "type", "t", "", typeUsage)
	insertCmd.Flags().StringVarP(&linkStatus, "linkStatus", "s", "", statusUsage)
	insertCmd.Flags().StringVarP(&externalChannelId, "externalChannelId", "e", "", extCidUsage)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("linkingToken")
	_ = insertCmd.MarkFlagRequired("type")
	cmd.AddMutationFlags(insertCmd)
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(c, "Would insert third-party link")
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
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}
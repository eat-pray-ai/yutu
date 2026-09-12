// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"encoding/json/jsontext"
	"fmt"
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
	insertConfirm = "Create third-party link: %s"
	insertShort   = "Create a third-party link"
	insertLong    = "Create a third-party link. Use this tool to create a third-party link."
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
			Default: jsontext.Value(`["snippet","status"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
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
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input thirdPartyLink.ThirdPartyLink) string {
					return fmt.Sprintf(insertConfirm, input.LinkingToken)
				},
				func(input thirdPartyLink.ThirdPartyLink, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	thirdPartyLinkCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&linkingToken, "linkingToken", "l", "", ltUsage)
	insertCmd.Flags().StringVarP(&linkType, "type", "t", "", typeUsage)
	insertCmd.Flags().StringVarP(&linkStatus, "linkStatus", "s", "", statusUsage)
	insertCmd.Flags().StringVarP(
		&externalChannelId, "externalChannelId", "e", "", extCidUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("linkingToken")
	_ = insertCmd.MarkFlagRequired("type")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(insertConfirm, linkingToken))
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
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

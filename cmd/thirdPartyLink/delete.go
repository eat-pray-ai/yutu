// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/thirdPartyLink"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool    = "thirdPartyLink-delete"
	deleteShort   = "Delete a third-party link"
	deleteLong    = "Delete a third-party link. Use this tool to delete a link between a YouTube channel and a third-party service."
	deleteExample = `# Delete a third-party link
yutu thirdPartyLink delete --linkingToken abc123 --type channelToStoreLink`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"linking_token", "type"},
	Properties: map[string]*jsonschema.Schema{
		"linking_token":       {Type: "string", Description: ltUsage},
		"type":                {Type: "string", Description: typeUsage, Enum: []any{"linkUnspecified", "channelToStoreLink"}},
		"external_channel_id": {Type: "string", Description: extCidUsage},
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
			deleteTool, func(input thirdPartyLink.ThirdPartyLink, writer io.Writer) error {
				return input.Delete(writer)
			},
		),
	)
	thirdPartyLinkCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&linkingToken, "linkingToken", "l", "", ltUsage)
	deleteCmd.Flags().StringVarP(&linkType, "type", "t", "", typeUsage)
	deleteCmd.Flags().StringVarP(&externalChannelId, "externalChannelId", "e", "", extCidUsage)

	_ = deleteCmd.MarkFlagRequired("linkingToken")
	_ = deleteCmd.MarkFlagRequired("type")
	cmd.AddMutationFlags(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Long:    deleteLong,
	Example: deleteExample,
	Run: func(c *cobra.Command, args []string) {
		err := cmd.Confirm(c, "Would delete third-party link: %s", linkingToken)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := thirdPartyLink.NewThirdPartyLink(
			thirdPartyLink.WithLinkingToken(linkingToken),
			thirdPartyLink.WithType(linkType),
			thirdPartyLink.WithExternalChannelId(externalChannelId),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}
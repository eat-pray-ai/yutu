// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
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
	deleteTool    = "thirdPartyLink-delete"
	deleteConfirm = "Delete third-party link(s): %s"
	deleteShort   = "Delete third-party links"
	deleteLong    = "Delete third-party links. Use this tool to delete third-party links."
	deleteExample = `  # Delete a third-party link
  yutu thirdPartyLink delete --linkingToken abc123 --type channelToStoreLink`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"linking_token", "type"},
	Properties: map[string]*jsonschema.Schema{
		"linking_token": {Type: "string", Description: ltUsage},
		"type": {
			Type: "string", Description: typeUsage,
			Enum: []any{"linkUnspecified", "channelToStoreLink"},
		},
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
		}, cobramcp.GenToolHandlerWithMRTR(
			deleteTool, cobramcp.ConfirmThen(
				func(input thirdPartyLink.ThirdPartyLink) string {
					return fmt.Sprintf(deleteConfirm, input.LinkingToken)
				},
				func(input thirdPartyLink.ThirdPartyLink, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	thirdPartyLinkCmd.AddCommand(newDeleteCmd())
}

type deleteOpts struct {
	linkingToken      string
	linkType          string
	externalChannelId string
	yes               bool
}

func (o *deleteOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.linkingToken, "linkingToken", "l", "", ltUsage)
	c.Flags().StringVarP(&o.linkType, "type", "t", "", typeUsage)
	c.Flags().StringVarP(
		&o.externalChannelId, "externalChannelId", "e", "", extCidUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("linkingToken")
	_ = c.MarkFlagRequired("type")
}

func (o *deleteOpts) build() *thirdPartyLink.ThirdPartyLink {
	return thirdPartyLink.NewThirdPartyLink(
		thirdPartyLink.WithLinkingToken(o.linkingToken),
		thirdPartyLink.WithType(o.linkType),
		thirdPartyLink.WithExternalChannelId(o.externalChannelId),
	).(*thirdPartyLink.ThirdPartyLink)
}

func newDeleteCmd() *cobra.Command {
	var opts deleteOpts

	c := &cobra.Command{
		Use:     "delete",
		Short:   deleteShort,
		Long:    deleteLong,
		Example: deleteExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(deleteConfirm, opts.linkingToken),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Delete(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

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
	updateTool    = "thirdPartyLink-update"
	updateConfirm = "Update third-party link: %s"
	updateShort   = "Update a third-party link"
	updateLong    = "Update a third-party link. Use this tool to update the status or type of a link between a YouTube channel and a third-party service."
	updateExample = `  # Update a third-party link status
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input thirdPartyLink.ThirdPartyLink) string {
					return fmt.Sprintf(updateConfirm, input.LinkingToken)
				},
				func(input thirdPartyLink.ThirdPartyLink, w io.Writer) error {
					return input.Update(w)
				},
			),
		),
	)
	thirdPartyLinkCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	linkingToken      string
	linkType          string
	linkStatus        string
	externalChannelId string
	parts             []string
	output            string
	yes               bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.linkingToken, "linkingToken", "l", "", ltUsage)
	c.Flags().StringVarP(&o.linkType, "type", "t", "", typeUsage)
	c.Flags().StringVarP(&o.linkStatus, "linkStatus", "s", "", statusUsage)
	c.Flags().StringVarP(
		&o.externalChannelId, "externalChannelId", "e", "", extCidUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("linkingToken")
}

func (o *updateOpts) build() *thirdPartyLink.ThirdPartyLink {
	return thirdPartyLink.NewThirdPartyLink(
		thirdPartyLink.WithLinkingToken(o.linkingToken),
		thirdPartyLink.WithType(o.linkType),
		thirdPartyLink.WithLinkStatus(o.linkStatus),
		thirdPartyLink.WithExternalChannelId(o.externalChannelId),
		thirdPartyLink.WithParts(o.parts),
		thirdPartyLink.WithOutput(o.output),
	).(*thirdPartyLink.ThirdPartyLink)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, opts.linkingToken),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

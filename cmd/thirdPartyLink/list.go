// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thirdPartyLink"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "thirdPartyLink-list"
	listShort   = "List third-party links"
	listLong    = "List third-party links. Use this tool to list links between a YouTube channel and a third-party service."
	listExample = `  # List all third-party links
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
		"linking_token":       {Type: "string", Description: ltUsage},
		"type":                {Type: "string", Description: typeUsage, Enum: []any{"linkUnspecified", "channelToStoreLink"}},
		"external_channel_id": {Type: "string", Description: extCidUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["snippet","status"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: jsontext.Value(`"json"`),
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
			listTool,
			func(input thirdPartyLink.ThirdPartyLink, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	thirdPartyLinkCmd.AddCommand(newListCmd())
}

type listOpts struct {
	linkingToken      string
	linkType          string
	externalChannelId string
	parts             []string
	output            string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.linkingToken, "linkingToken", "l", "", ltUsage)
	c.Flags().StringVarP(&o.linkType, "type", "t", "", typeUsage)
	c.Flags().StringVarP(
		&o.externalChannelId, "externalChannelId", "e", "", extCidUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *thirdPartyLink.ThirdPartyLink {
	return thirdPartyLink.NewThirdPartyLink(
		thirdPartyLink.WithLinkingToken(o.linkingToken),
		thirdPartyLink.WithType(o.linkType),
		thirdPartyLink.WithExternalChannelId(o.externalChannelId),
		thirdPartyLink.WithParts(o.parts),
		thirdPartyLink.WithOutput(o.output),
	).(*thirdPartyLink.ThirdPartyLink)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

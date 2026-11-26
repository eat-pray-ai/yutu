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
	insertExample = `  # Insert a new third-party link
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
	thirdPartyLinkCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	linkingToken      string
	linkType          string
	linkStatus        string
	externalChannelId string
	parts             []string
	output            string
	yes               bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
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
	_ = c.MarkFlagRequired("type")
}

func (o *insertOpts) build() *thirdPartyLink.ThirdPartyLink {
	return thirdPartyLink.NewThirdPartyLink(
		thirdPartyLink.WithLinkingToken(o.linkingToken),
		thirdPartyLink.WithType(o.linkType),
		thirdPartyLink.WithLinkStatus(o.linkStatus),
		thirdPartyLink.WithExternalChannelId(o.externalChannelId),
		thirdPartyLink.WithParts(o.parts),
		thirdPartyLink.WithOutput(o.output),
	).(*thirdPartyLink.ThirdPartyLink)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.linkingToken),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	listTool     = "channelSection-list"
	listIdsUsage = "Return the channel sections with the given ids"
	listShort    = "List channel sections"
	listLong     = "List channel sections. Use this tool to list channel sections."
	listExample  = `  # List my channel sections
  yutu channelSection list --mine
  # List channel sections by channel ID
  yutu channelSection list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
  # List specific channel sections in JSON format
  yutu channelSection list --ids abc123,def456 --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"channel_id":                 {Type: "string", Description: cidUsage},
		"hl":                         {Type: "string", Description: hlUsage},
		"mine":                       {Type: "boolean", Description: mineUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet"]`),
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
			func(input channelSection.ChannelSection, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	channelSectionCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids       []string
	channelId string
	hl        string
	mine      bool
	oboco     string
	parts     []string
	output    string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().BoolVarP(&o.mine, "mine", "M", false, mineUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build(f *pflag.FlagSet) *channelSection.ChannelSection {
	return channelSection.NewChannelSection(
		channelSection.WithIds(o.ids),
		channelSection.WithChannelId(o.channelId),
		channelSection.WithHl(o.hl),
		utils.OptIfChanged(f, "mine", &o.mine, channelSection.WithMine),
		channelSection.WithOBOCO(o.oboco),
		channelSection.WithParts(o.parts),
		channelSection.WithOutput(o.output),
	).(*channelSection.ChannelSection)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

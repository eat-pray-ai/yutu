// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "channel-list"
	listIdsUsage = "Return the channels with the specified Ids"
	listShort    = "List channels"
	listLong     = "List channels. Use this tool to list channels."
	listExample  = `  # List my channel info
  yutu channel list --for mine
  # List channel by handle
  yutu channel list --forHandle @GoogleDevelopers
  # List channel by ID in JSON format
  yutu channel list --ids UC_x5XG1OV2P6uZZ5FSM9Ttw --output json
  # List channels in table format
  yutu channel list --maxResults 10 --output table`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"category_id":  {Type: "string", Description: cidUsage},
		"for_handle":   {Type: "string", Description: fhUsage},
		"for_username": {Type: "string", Description: fuUsage},
		"hl":           {Type: "string", Description: hlUsage},
		"ids": {
			Type: "array", Items: &jsonschema.Schema{Type: "string"},
			Description: listIdsUsage,
			Default:     jsontext.Value(`[]`),
		},
		"for": {
			Type: "string", Description: forUsage,
			Enum: []any{"managedByMe", "mine", "mySubscribers"},
		},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","status"]`),
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
			listTool, func(input channel.Channel, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	channelCmd.AddCommand(newListCmd())
}

type listOpts struct {
	categoryId  string
	forHandle   string
	forUsername string
	hl          string
	ids         []string
	channelFor  string
	maxResults  int64
	oboco       string
	parts       []string
	output      string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.categoryId, "categoryId", "g", "", cidUsage)
	c.Flags().StringVarP(&o.forHandle, "forHandle", "d", "", fhUsage)
	c.Flags().StringVarP(&o.forUsername, "forUsername", "u", "", fuUsage)
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.channelFor, "for", "F", "", forUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *channel.Channel {
	return channel.NewChannel(
		channel.WithCategoryId(o.categoryId),
		channel.WithForHandle(o.forHandle),
		channel.WithForUsername(o.forUsername),
		channel.WithHl(o.hl),
		channel.WithIds(o.ids),
		channel.WithFor(o.channelFor),
		channel.WithMaxResults(o.maxResults),
		channel.WithOBOCO(o.oboco),
		channel.WithParts(o.parts),
		channel.WithOutput(o.output),
	).(*channel.Channel)
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "subscription-list"
	listIdsUsage = "Return the subscriptions with the given ids for Stubby or Apiary"
	listCidUsage = "Return the subscriptions of the given channel owner"
	listShort    = "List subscriptions"
	listLong     = "List subscriptions. Use this tool to list subscriptions."
	listExample  = `  # List my subscriptions
  yutu subscription list --for mine
  # List subscriptions by channel ID with limit
  yutu subscription list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10
  # List specific subscriptions in JSON format
  yutu subscription list --ids abc123,def456 --output json
  # List subscriptions for a channel alphabetically
  yutu subscription list --forChannelId UC_x5XG1OV2P6uZZ5FSM9Ttw --order alphabetical`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"channel_id":     {Type: "string", Description: listCidUsage},
		"for_channel_id": {Type: "string", Description: fcidUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"for": {
			Type: "string", Description: forUsage,
			Enum: []any{"mine", "myRecentSubscribers", "mySubscribers"},
		},
		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
		},
		"order": {
			Type: "string", Description: orderUsage,
			Enum: []any{
				"subscriptionOrderUnspecified", "relevance", "unread", "alphabetical",
			},
			Default: jsontext.Value(`"relevance"`),
		},
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
			listTool, func(input subscription.Subscription, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	subscriptionCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                           []string
	channelId                     string
	forChannelId                  string
	maxResults                    int64
	subscriptionFor               string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	order                         string
	parts                         []string
	output                        string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", listCidUsage)
	c.Flags().StringVarP(&o.forChannelId, "forChannelId", "C", "", fcidUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(&o.subscriptionFor, "for", "F", "", forUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		pkg.OBOCOCUsage,
	)
	c.Flags().StringVarP(&o.order, "order", "O", "relevance", orderUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *subscription.Subscription {
	return subscription.NewSubscription(
		subscription.WithIds(o.ids),
		subscription.WithChannelId(o.channelId),
		subscription.WithForChannelId(o.forChannelId),
		subscription.WithMaxResults(o.maxResults),
		subscription.WithFor(o.subscriptionFor),
		subscription.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		subscription.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		subscription.WithOrder(o.order),
		subscription.WithParts(o.parts),
		subscription.WithOutput(o.output),
	).(*subscription.Subscription)
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "subscription-insert"
	insertCidUsage = "ID of the channel to subscribe to"
	insertConfirm  = "Subscribe to channel: %s"
	insertShort    = "Subscribe to a channel"
	insertLong     = "Subscribe to a channel. Use this tool to subscribe to a YouTube channel."
	insertExample  = `  # Subscribe to a channel
  yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw
  # Subscribe with a title
  yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5X --title 'Google Developers'`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"subscriber_channel_id", "channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"subscriber_channel_id": {Type: "string", Description: scidUsage},
		"description":           {Type: "string", Description: descUsage},
		"channel_id":            {Type: "string", Description: insertCidUsage},
		"title":                 {Type: "string", Description: titleUsage},
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
				func(input subscription.Subscription) string {
					return fmt.Sprintf(insertConfirm, input.ChannelId)
				},
				func(input subscription.Subscription, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	subscriptionCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	subscriberChannelId string
	description         string
	channelId           string
	title               string
	output              string
	yes                 bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&o.subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", insertCidUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("subscriberChannelId")
	_ = c.MarkFlagRequired("channelId")
}

func (o *insertOpts) build() *subscription.Subscription {
	return subscription.NewSubscription(
		subscription.WithSubscriberChannelId(o.subscriberChannelId),
		subscription.WithDescription(o.description),
		subscription.WithChannelId(o.channelId),
		subscription.WithTitle(o.title),
		subscription.WithOutput(o.output),
	).(*subscription.Subscription)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.channelId))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

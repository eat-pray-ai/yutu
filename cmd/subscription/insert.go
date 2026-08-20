// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"encoding/json"
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
	insertCidUsage = "ID of the channel to be subscribed"
	insertShort    = "Insert a new subscription"
	insertLong     = "Insert a new subscription. Use this tool to insert a new subscription."
	insertExample  = `# Subscribe to a channel
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
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
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
		}, cobramcp.GenToolHandler(
			insertTool, func(input subscription.Subscription, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Insert(writer)
			},
		),
	)
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf("Would insert subscription to channel: %s", channelId)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := subscription.NewSubscription(
			subscription.WithSubscriberChannelId(subscriberChannelId),
			subscription.WithDescription(description),
			subscription.WithChannelId(channelId),
			subscription.WithTitle(title),
			subscription.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

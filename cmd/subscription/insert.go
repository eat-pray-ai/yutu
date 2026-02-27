// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"encoding/json"
	"io"

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
	insertShort    = "Insert a YouTube subscription"
	insertLong     = "Insert a YouTube subscription\n\nExamples:\n  yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw\n  yutu subscription insert --subscriberChannelId UC_abc --channelId UC_x5X --title 'Google Developers'"
	insertCidUsage = "ID of the channel to be subscribed"
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
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
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
		}, cmd.GenToolHandler(
			insertTool, func(input subscription.Subscription, writer io.Writer) error {
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
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := subscription.NewSubscription(
			subscription.WithSubscriberChannelId(subscriberChannelId),
			subscription.WithDescription(description),
			subscription.WithChannelId(channelId),
			subscription.WithTitle(title),
			subscription.WithOutput(output),
			subscription.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Insert(cmd.OutOrStdout()), cmd)
	},
}

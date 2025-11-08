// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "subscription-insert"
	insertShort    = "Insert a YouTube subscription"
	insertLong     = "Insert a YouTube subscription"
	insertCidUsage = "ID of the channel to be subscribed"
)

type insertIn struct {
	SubscriberChannelId string `json:"subscriberChannelId"`
	Description         string `json:"description"`
	ChannelId           string `json:"channelId"`
	Title               string `json:"title"`
	Output              string `json:"output"`
	Jsonpath            string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"subscriberChannelId", "description", "channelId", "title",
		"output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"subscriberChannelId": {
			Type: "string", Description: scidUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"channelId": {
			Type: "string", Description: insertCidUsage,
			Default: json.RawMessage(`""`),
		},
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	subscriberChannelId = input.SubscriberChannelId
	description = input.Description
	channelId = input.ChannelId
	title = input.Title
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithSubscriberChannelId(subscriberChannelId),
		subscription.WithDescription(description),
		subscription.WithChannelId(channelId),
		subscription.WithTitle(title),
		subscription.WithService(nil),
	)

	return s.Insert(output, jpath, writer)
}

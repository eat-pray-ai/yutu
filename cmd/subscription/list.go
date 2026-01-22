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
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "subscription-list"
	listShort    = "List subscriptions' info"
	listLong     = "List subscriptions' info, such as id, title, etc"
	listIdsUsage = "Return the subscriptions with the given ids for Stubby or Apiary"
	listCidUsage = "Return the subscriptions of the given channel owner"
)

type listIn struct {
	Ids                           []string `json:"ids"`
	ChannelId                     string   `json:"channelId"`
	ForChannelId                  string   `json:"forChannelId"`
	MaxResults                    int64    `json:"maxResults"`
	Mine                          *string  `json:"mine,omitempty"`
	MyRecentSubscribers           *string  `json:"myRecentSubscribers,omitempty"`
	MySubscribers                 *string  `json:"mySubscribers,omitempty"`
	OnBehalfOfContentOwner        string   `json:"onBehalfOfContentOwner"`
	OnBehalfOfContentOwnerChannel string   `json:"onBehalfOfContentOwnerChannel"`
	Order                         string   `json:"order"`
	Parts                         []string `json:"parts"`
	Output                        string   `json:"output"`
	Jsonpath                      string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"channelId": {
			Type: "string", Description: listCidUsage,
			Default: json.RawMessage(`""`),
		},
		"forChannelId": {
			Type: "string", Description: fcidUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"mine": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: mineUsage, Default: json.RawMessage(`""`),
		},
		"myRecentSubscribers": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: mrsUsage, Default: json.RawMessage(`""`),
		},
		"mySubscribers": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: msUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwnerChannel": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"order": {
			Type: "string", Enum: []any{
				"subscriptionOrderUnspecified", "relevance", "unread", "alphabetical",
			},
			Description: orderUsage, Default: json.RawMessage(`"relevance"`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
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
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
	)
	subscriptionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&forChannelId, "forChannelId", "C", "", fcidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		myRecentSubscribers, "myRecentSubscribers", "R", false, mrsUsage,
	)
	listCmd.Flags().BoolVarP(mySubscribers, "mySubscribers", "S", false, msUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "relevance", orderUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func listHandler(
	ctx context.Context, req *mcp.CallToolRequest, input listIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: listTool, MinInterval: time.Second},
		),
	)

	ids = input.Ids
	channelId = input.ChannelId
	forChannelId = input.ForChannelId
	maxResults = input.MaxResults
	mine = utils.BoolPtr(*input.Mine)
	myRecentSubscribers = utils.BoolPtr(*input.MyRecentSubscribers)
	mySubscribers = utils.BoolPtr(*input.MySubscribers)
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	onBehalfOfContentOwnerChannel = input.OnBehalfOfContentOwnerChannel
	order = input.Order
	parts = input.Parts
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithIDs(ids),
		subscription.WithChannelId(channelId),
		subscription.WithForChannelId(forChannelId),
		subscription.WithMaxResults(maxResults),
		subscription.WithMine(mine),
		subscription.WithMyRecentSubscribers(myRecentSubscribers),
		subscription.WithMySubscribers(mySubscribers),
		subscription.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		subscription.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		subscription.WithOrder(order),
		subscription.WithService(nil),
	)

	return s.List(parts, output, jsonpath, writer)
}

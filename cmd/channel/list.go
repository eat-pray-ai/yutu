// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "channel-list"
	listShort    = "List channel's info"
	listLong     = "List channel's info, such as title, description, etc."
	listIdsUsage = "Return the channels with the specified IDs"
)

type listIn struct {
	CategoryId             string   `json:"categoryId"`
	ForHandle              string   `json:"forHandle"`
	ForUsername            string   `json:"forUsername"`
	Hl                     string   `json:"hl"`
	Ids                    []string `json:"ids"`
	ManagedByMe            *string  `json:"managedByMe,omitempty"`
	MaxResults             int64    `json:"maxResults"`
	Mine                   *string  `json:"mine,omitempty"`
	MySubscribers          *string  `json:"mySubscribers,omitempty"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Parts                  []string `json:"parts"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"categoryId": {
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"forHandle": {
			Type: "string", Description: fhUsage,
			Default: json.RawMessage(`""`),
		},
		"forUsername": {
			Type: "string", Description: fuUsage,
			Default: json.RawMessage(`""`),
		},
		"hl": {
			Type: "string", Description: hlUsage,
			Default: json.RawMessage(`""`),
		},
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"managedByMe": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: mbmUsage, Default: json.RawMessage(`""`),
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
		"mySubscribers": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: msUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet","status"]`),
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
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&categoryId, "categoryId", "g", "", cidUsage,
	)
	listCmd.Flags().StringVarP(
		&forHandle, "forHandle", "d", "", fhUsage,
	)
	listCmd.Flags().StringVarP(
		&forUsername, "forUsername", "u", "", fuUsage,
	)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().BoolVarP(
		managedByMe, "managedByMe", "E", false, mbmUsage,
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, pkg.MRUsage,
	)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		mySubscribers, "mySubscribers", "S", false, msUsage,
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
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

	categoryId = input.CategoryId
	forHandle = input.ForHandle
	forUsername = input.ForUsername
	hl = input.Hl
	ids = input.Ids
	managedByMe = utils.BoolPtr(*input.ManagedByMe)
	maxResults = input.MaxResults
	mine = utils.BoolPtr(*input.Mine)
	mySubscribers = utils.BoolPtr(*input.MySubscribers)
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
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
	c := channel.NewChannel(
		channel.WithCategoryId(categoryId),
		channel.WithForHandle(forHandle),
		channel.WithForUsername(forUsername),
		channel.WithHl(hl),
		channel.WithIDs(ids),
		channel.WithChannelManagedByMe(managedByMe),
		channel.WithMaxResults(maxResults),
		channel.WithMine(mine),
		channel.WithMySubscribers(mySubscribers),
		channel.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channel.WithService(nil),
	)

	return c.List(parts, output, jsonpath, writer)
}

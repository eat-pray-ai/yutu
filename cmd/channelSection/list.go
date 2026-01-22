// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "channelSection-list"
	listShort    = "List channel sections"
	listLong     = "List channel sections"
	listIdsUsage = "Return the channel sections with the given ids"
)

type listIn struct {
	Ids                    []string `json:"ids"`
	ChannelId              string   `json:"channelId"`
	Hl                     string   `json:"hl"`
	Mine                   *string  `json:"mine,omitempty"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Parts                  []string `json:"parts"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
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
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"hl": {
			Type: "string", Description: hlUsage,
			Default: json.RawMessage(`""`),
		},
		"mine": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: mineUsage, Default: json.RawMessage(`""`),
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
	channelSectionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", false, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
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
	hl = input.Hl
	mine = utils.BoolPtr(*input.Mine)
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
	cs := channelSection.NewChannelSection(
		channelSection.WithIDs(ids),
		channelSection.WithChannelId(channelId),
		channelSection.WithHl(hl),
		channelSection.WithMine(mine),
		channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channelSection.WithService(nil),
	)

	return cs.List(parts, output, jsonpath, writer)
}

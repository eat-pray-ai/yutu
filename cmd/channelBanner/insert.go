// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channelBanner"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool = "channelBanner-insert"
)

type insertIn struct {
	ChannelId                     string `json:"channelId"`
	File                          string `json:"file"`
	OnBehalfOfContentOwner        string `json:"onBehalfOfContentOwner"`
	OnBehalfOfContentOwnerChannel string `json:"onBehalfOfContentOwnerChannel"`
	Output                        string `json:"output"`
	Jsonpath                      string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channelId", "file"},
	Properties: map[string]*jsonschema.Schema{
		"channelId": {
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"file": {
			Type: "string", Description: fileUsage,
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwnerChannel": {
			Type: "string", Description: "",
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
			Name: insertTool, Title: short, Description: long,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	channelBannerCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("file")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: short,
	Long:  long,
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

	channelId = input.ChannelId
	file = input.File
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	onBehalfOfContentOwnerChannel = input.OnBehalfOfContentOwnerChannel
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	cb := channelBanner.NewChannelBanner(
		channelBanner.WithChannelId(channelId),
		channelBanner.WithFile(file),
		channelBanner.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channelBanner.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		channelBanner.WithService(nil),
	)

	return cb.Insert(output, jsonpath, writer)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	unsetShort = "Unset watermark for channel's video"
	unsetLong  = "Unset watermark for channel's video by channel id"
)

type unsetIn struct {
	ChannelId string `json:"channelId"`
}

var unsetInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channelId"},
	Properties: map[string]*jsonschema.Schema{
		"channelId": {
			Type:        "string",
			Description: cidUsage,
			Default:     json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name:        "watermark-unset",
			Title:       unsetShort,
			Description: unsetLong,
			InputSchema: unsetInSchema,
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, unsetHandler,
	)
	watermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	_ = unsetCmd.MarkFlagRequired("channelId")
}

var unsetCmd = &cobra.Command{
	Use:   "unset",
	Short: unsetShort,
	Long:  unsetLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := unset(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func unsetHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input unsetIn,
) (*mcp.CallToolResult, any, error) {
	channelId = input.ChannelId

	slog.InfoContext(ctx, "watermark unset started")

	var writer bytes.Buffer
	err := unset(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "watermark unset failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "watermark unset completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func unset(writer io.Writer) error {
	w := watermark.NewWatermark(
		watermark.WithChannelId(channelId),
		watermark.WithService(nil),
	)

	return w.Unset(writer)
}

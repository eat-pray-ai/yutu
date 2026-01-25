// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"bytes"
	"context"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	setTool  = "watermark-set"
	setShort = "Set watermark for channel's video"
	setLong  = "Set watermark for channel's video by channel id"
)

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channel_id", "file"},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
		"file":       {Type: "string", Description: fileUsage},
		"in_video_position": {
			Type: "string", Description: ivpUsage,
			Enum: []any{"topLeft", "topRight", "bottomLeft", "bottomRight", ""},
		},
		"duration_ms": {Type: "number", Description: dmUsage},
		"offset_ms":   {Type: "number", Description: omUsage},
		"offset_type": {
			Type: "string", Description: otUsage,
			Enum: []any{"offsetFromStart", "offsetFromEnd", ""},
		},
		"on_behalf_of_content_owner": {Type: "string"},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name:        setTool,
			Title:       setShort,
			Description: setLong,
			InputSchema: setInSchema,
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, setHandler,
	)
	watermarkCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(
		&inVideoPosition, "inVideoPosition", "p", "", ivpUsage,
	)
	setCmd.Flags().Uint64VarP(&durationMs, "durationMs", "d", 0, dmUsage)
	setCmd.Flags().Uint64VarP(&offsetMs, "offsetMs", "m", 0, omUsage)
	setCmd.Flags().StringVarP(&offsetType, "offsetType", "t", "", otUsage)
	setCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = setCmd.MarkFlagRequired("channelId")
	_ = setCmd.MarkFlagRequired("file")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: setShort,
	Long:  setLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := watermark.NewWatermark(
			watermark.WithChannelId(channelId),
			watermark.WithFile(file),
			watermark.WithInVideoPosition(inVideoPosition),
			watermark.WithDurationMs(durationMs),
			watermark.WithOffsetMs(offsetMs),
			watermark.WithOffsetType(offsetType),
			watermark.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		err := input.Set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setHandler(
	ctx context.Context, req *mcp.CallToolRequest, input watermark.Watermark,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: setTool, MinInterval: time.Second},
		),
	)

	var writer bytes.Buffer
	err := input.Set(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

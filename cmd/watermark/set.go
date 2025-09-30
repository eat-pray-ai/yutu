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
	setShort = "Set watermark for channel's video"
	setLong  = "Set watermark for channel's video by channel id"
)

type setIn struct {
	ChannelId              string `json:"channelId"`
	File                   string `json:"file"`
	InVideoPosition        string `json:"inVideoPosition"`
	DurationMs             int64  `json:"durationMs"`
	OffsetMs               int64  `json:"offsetMs"`
	OffsetType             string `json:"offsetType"`
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"`
}

var setInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"channelId", "file", "inVideoPosition", "durationMs",
		"offsetMs", "offsetType", "onBehalfOfContentOwner",
	},
	Properties: map[string]*jsonschema.Schema{
		"channelId": {
			Type:        "string",
			Description: cidUsage,
			Default:     json.RawMessage(`""`),
		},
		"file": {
			Type:        "string",
			Description: fileUsage,
			Default:     json.RawMessage(`""`),
		},
		"inVideoPosition": {
			Type:        "string",
			Enum:        []any{"topLeft", "topRight", "bottomLeft", "bottomRight", ""},
			Description: ivpUsage,
			Default:     json.RawMessage(`""`),
		},
		"durationMs": {
			Type:        "number",
			Description: dmUsage,
			Default:     json.RawMessage("0"),
		},
		"offsetMs": {
			Type:        "number",
			Description: omUsage,
			Default:     json.RawMessage("0"),
		},
		"offsetType": {
			Type:        "string",
			Enum:        []any{"offsetFromStart", "offsetFromEnd", ""},
			Description: otUsage,
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name:        "watermark-set",
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
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input setIn,
) (*mcp.CallToolResult, any, error) {
	channelId = input.ChannelId
	file = input.File
	inVideoPosition = input.InVideoPosition
	durationMs = uint64(input.DurationMs)
	offsetMs = uint64(input.OffsetMs)
	offsetType = input.OffsetType
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	slog.InfoContext(ctx, "watermark set started")

	var writer bytes.Buffer
	err := set(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "watermark set failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "watermark set completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func set(writer io.Writer) error {
	w := watermark.NewWatermark(
		watermark.WithChannelId(channelId),
		watermark.WithFile(file),
		watermark.WithInVideoPosition(inVideoPosition),
		watermark.WithDurationMs(durationMs),
		watermark.WithOffsetMs(offsetMs),
		watermark.WithOffsetType(offsetType),
		watermark.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		watermark.WithService(nil),
	)

	return w.Set(writer)
}

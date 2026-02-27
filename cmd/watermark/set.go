// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	setTool  = "watermark-set"
	setShort = "Set watermark for channel's video"
	setLong  = "Set watermark for channel's video by channel id\n\nExamples:\n  yutu watermark set --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file watermark.png\n  yutu watermark set --channelId UC_x5X --file logo.png --inVideoPosition bottomRight --offsetType offsetFromEnd --offsetMs 1000\n  yutu watermark set --channelId UC_x5X --file logo.png --durationMs 5000 --offsetMs 2000"
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
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			setTool, func(input watermark.Watermark, writer io.Writer) error {
				return input.Set(writer)
			},
		),
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
		utils.HandleCmdError(input.Set(cmd.OutOrStdout()), cmd)
	},
}

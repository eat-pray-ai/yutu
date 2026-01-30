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
	unsetTool  = "watermark-unset"
	unsetShort = "Unset watermark for channel's video"
	unsetLong  = "Unset watermark for channel's video by channel id"
)

var unsetInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: unsetTool, Title: unsetShort, Description: unsetLong,
			InputSchema: unsetInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			unsetTool, func(input watermark.Watermark, writer io.Writer) error {
				return input.Unset(writer)
			},
		),
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
		input := watermark.NewWatermark(watermark.WithChannelId(channelId))
		utils.HandleCmdError(input.Unset(cmd.OutOrStdout()), cmd)
	},
}

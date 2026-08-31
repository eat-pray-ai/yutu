// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	unsetTool    = "watermark-unset"
	unsetConfirm = "Unset watermark for channel: %s"
	unsetShort   = "Unset a watermark for channel's videos"
	unsetLong    = "Unset a watermark for channel's videos. Use this tool to unset a watermark for a channel's videos."
	unsetExample = `# Unset watermark for a channel
yutu watermark unset --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw`
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
				DestructiveHint: new(true),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			unsetTool, cobramcp.ConfirmThen(
				func(input watermark.Watermark) string {
					return fmt.Sprintf(unsetConfirm, input.ChannelId)
				},
				func(input watermark.Watermark, w io.Writer) error {
					return input.Unset(w)
				},
			),
		),
	)
	watermarkCmd.AddCommand(unsetCmd)

	unsetCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	unsetCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = unsetCmd.MarkFlagRequired("channelId")
}

var unsetCmd = &cobra.Command{
	Use:     "unset",
	Short:   unsetShort,
	Long:    unsetLong,
	Example: unsetExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(unsetConfirm, channelId))
	},
	Run: func(c *cobra.Command, _ []string) {
		input := watermark.NewWatermark(watermark.WithChannelId(channelId))
		utils.HandleCmdError(input.Unset(c.OutOrStdout()), c)
	},
}

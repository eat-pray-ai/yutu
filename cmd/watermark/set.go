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
	setTool    = "watermark-set"
	setConfirm = "Set watermark %s for channel %s"
	setShort   = "Set a channel watermark"
	setLong    = "Set a channel watermark. Use this tool to set a channel watermark."
	setExample = `  # Set a watermark for a channel
  yutu watermark set --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file watermark.png
  # Set a watermark with position and offset
  yutu watermark set --channelId UC_x5X --file logo.png --inVideoPosition bottomRight --offsetType offsetFromEnd --offsetMs 1000
  # Set a watermark with duration
  yutu watermark set --channelId UC_x5X --file logo.png --durationMs 5000 --offsetMs 2000`
)

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channel_id", "file"},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
		"file":       {Type: "string", Description: fileUsage},
		"in_video_position": {
			Type: "string", Description: ivpUsage,
			Enum: []any{"topLeft", "topRight", "bottomLeft", "bottomRight"},
		},
		"duration_ms": {Type: "number", Description: dmUsage, Minimum: new(float64(0))},
		"offset_ms":   {Type: "number", Description: omUsage, Minimum: new(float64(0))},
		"offset_type": {
			Type: "string", Description: otUsage,
			Enum: []any{"offsetFromStart", "offsetFromEnd"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: setTool, Title: setShort, Description: setLong,
			InputSchema: setInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			setTool, cobramcp.ConfirmThen(
				func(input watermark.Watermark) string {
					return fmt.Sprintf(setConfirm, input.File, input.ChannelId)
				},
				func(input watermark.Watermark, w io.Writer) error {
					return input.Set(w)
				},
			),
		),
	)
	watermarkCmd.AddCommand(newSetCmd())
}

type setOpts struct {
	channelId              string
	file                   string
	inVideoPosition        string
	durationMs             uint64
	offsetMs               uint64
	offsetType             string
	onBehalfOfContentOwner string
	yes                    bool
}

func (o *setOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(
		&o.inVideoPosition, "inVideoPosition", "p", "", ivpUsage,
	)
	c.Flags().Uint64VarP(&o.durationMs, "durationMs", "d", 0, dmUsage)
	c.Flags().Uint64VarP(&o.offsetMs, "offsetMs", "m", 0, omUsage)
	c.Flags().StringVarP(&o.offsetType, "offsetType", "t", "", otUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("channelId")
	_ = c.MarkFlagRequired("file")
}

func (o *setOpts) build() *watermark.Watermark {
	return watermark.NewWatermark(
		watermark.WithChannelId(o.channelId),
		watermark.WithFile(o.file),
		watermark.WithInVideoPosition(o.inVideoPosition),
		watermark.WithDurationMs(o.durationMs),
		watermark.WithOffsetMs(o.offsetMs),
		watermark.WithOffsetType(o.offsetType),
		watermark.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
	).(*watermark.Watermark)
}

func newSetCmd() *cobra.Command {
	var opts setOpts

	c := &cobra.Command{
		Use:     "set",
		Short:   setShort,
		Long:    setLong,
		Example: setExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(setConfirm, opts.file, opts.channelId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Set(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

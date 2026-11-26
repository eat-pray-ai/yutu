// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveStream"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	updateTool    = "liveStream-update"
	updateConfirm = "Update live stream: %s"
	updateIdUsage = "ID of the live stream to update"
	updateShort   = "Update a live stream"
	updateLong    = "Update a live stream. Use this tool to update an existing live stream's settings."
	updateExample = `  # Update live stream title
  yutu liveStream update --id stream123 --title "New Title"
  # Update live stream CDN settings
  yutu liveStream update --id stream123 --resolution 1080p --frameRate 60fps`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"frame_rate": {
			Type: "string", Description: frUsage,
			Enum: []any{"30fps", "60fps", "variable"},
		},
		"ingestion_type": {
			Type: "string", Description: itUsage,
			Enum: []any{"rtmp", "dash", "webrtc", "hls"},
		},
		"resolution": {
			Type: "string", Description: resUsage,
			Enum: []any{"240p", "360p", "480p", "720p", "1080p", "1440p", "2160p", "variable"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input liveStream.LiveStream) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input liveStream.LiveStream, w io.Writer) error {
					liveStream.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	liveStreamCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids                           []string
	title                         string
	description                   string
	frameRate                     string
	ingestionType                 string
	resolution                    string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	output                        string
	yes                           bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(&o.frameRate, "frameRate", "f", "", frUsage)
	c.Flags().StringVarP(&o.ingestionType, "ingestionType", "I", "", itUsage)
	c.Flags().StringVarP(&o.resolution, "resolution", "r", "", resUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *liveStream.LiveStream {
	return liveStream.NewLiveStream(
		liveStream.WithIds(o.ids),
		liveStream.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, liveStream.WithDescription,
		),
		liveStream.WithFrameRate(o.frameRate),
		liveStream.WithIngestionType(o.ingestionType),
		liveStream.WithResolution(o.resolution),
		liveStream.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		liveStream.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		liveStream.WithMaxResults(1),
		liveStream.WithOutput(o.output),
	).(*liveStream.LiveStream)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

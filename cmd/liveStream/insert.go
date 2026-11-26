// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

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
	insertTool    = "liveStream-insert"
	insertConfirm = "Create live stream: %s"
	insertShort   = "Create a live stream"
	insertLong    = "Create a live stream. Use this tool to create a live stream."
	insertExample = `  # Create a live stream with RTMP at 1080p
  yutu liveStream insert --title "My Stream" --ingestionType rtmp --resolution 1080p --frameRate 60fps
  # Create a live stream with description
  yutu liveStream insert --title "Gaming Stream" --description "Live gaming session" --ingestionType rtmp --resolution 720p`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"title", "ingestion_type", "resolution"},
	Properties: map[string]*jsonschema.Schema{
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
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["snippet","cdn","status"]`),
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input liveStream.LiveStream) string {
					return fmt.Sprintf(insertConfirm, input.Title)
				},
				func(input liveStream.LiveStream, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	liveStreamCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	title         string
	description   string
	frameRate     string
	ingestionType string
	resolution    string
	oboco         string
	obococ        string
	parts         []string
	output        string
	yes           bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(&o.frameRate, "frameRate", "f", "", frUsage)
	c.Flags().StringVarP(&o.ingestionType, "ingestionType", "I", "", itUsage)
	c.Flags().StringVarP(&o.resolution, "resolution", "r", "", resUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet", "cdn", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("ingestionType")
	_ = c.MarkFlagRequired("resolution")
}

func (o *insertOpts) build(f *pflag.FlagSet) *liveStream.LiveStream {
	return liveStream.NewLiveStream(
		liveStream.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, liveStream.WithDescription,
		),
		liveStream.WithFrameRate(o.frameRate),
		liveStream.WithIngestionType(o.ingestionType),
		liveStream.WithResolution(o.resolution),
		liveStream.WithOBOCO(o.oboco),
		liveStream.WithOBOCOC(o.obococ),
		liveStream.WithParts(o.parts),
		liveStream.WithOutput(o.output),
	).(*liveStream.LiveStream)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.title))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

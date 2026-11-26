// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveBroadcast"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	bindTool    = "liveBroadcast-bind"
	bindConfirm = "Bind live broadcast(s) %s to stream %s"
	bindShort   = "Bind a live broadcast to a stream"
	bindLong    = "Bind a live broadcast to a stream. Use this tool to bind or unbind a live stream to/from a live broadcast."
	bindExample = `  # Bind a broadcast to a stream
  yutu liveBroadcast bind --ids broadcast123 --streamId stream456
  # Unbind a broadcast (omit streamId)
  yutu liveBroadcast bind --ids broadcast123`
)

var bindInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live broadcasts to bind",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"stream_id":                  {Type: "string", Description: sidUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","contentDetails","status"]`),
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
			Name: bindTool, Title: bindShort, Description: bindLong,
			InputSchema: bindInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			bindTool, cobramcp.ConfirmThen(
				func(input liveBroadcast.LiveBroadcast) string {
					return fmt.Sprintf(
						bindConfirm, strings.Join(input.Ids, ", "), input.StreamId,
					)
				},
				func(input liveBroadcast.LiveBroadcast, w io.Writer) error {
					return input.Bind(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(newBindCmd())
}

type bindOpts struct {
	ids                           []string
	streamId                      string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	parts                         []string
	output                        string
	yes                           bool
}

func (o *bindOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(
		&o.ids, "ids", "i", nil, "IDs of the broadcasts to bind",
	)
	c.Flags().StringVarP(&o.streamId, "streamId", "s", "", sidUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "contentDetails", "status"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *bindOpts) build() *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithIds(o.ids),
		liveBroadcast.WithStreamId(o.streamId),
		liveBroadcast.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		liveBroadcast.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		liveBroadcast.WithParts(o.parts),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
}

func newBindCmd() *cobra.Command {
	var opts bindOpts

	c := &cobra.Command{
		Use:     "bind",
		Short:   bindShort,
		Long:    bindLong,
		Example: bindExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(bindConfirm, strings.Join(opts.ids, ", "), opts.streamId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Bind(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

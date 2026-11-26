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
	insertCuepointTool    = "liveBroadcast-insertCuepoint"
	insertCuepointConfirm = "Insert cuepoint into live broadcast(s) %s"
	insertCuepointShort   = "Insert a cuepoint into a live broadcast"
	insertCuepointLong    = "Insert a cuepoint into a live broadcast. Use this tool to insert an ad break cuepoint into a currently live broadcast."
	insertCuepointExample = `  # Insert an ad cuepoint of 30 seconds
  yutu liveBroadcast insertCuepoint --ids broadcast123 --cueType cueTypeAd --cueDurationSecs 30
  # Insert a cuepoint with offset
  yutu liveBroadcast insertCuepoint --ids broadcast123 --cueType cueTypeAd --cueInsertionOffsetMs 5000`
)

var insertCuepointInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live broadcasts to insert cuepoints into",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"cue_type": {
			Type: "string", Enum: []any{"cueTypeAd"},
			Description: ctUsage,
		},
		"cue_duration_secs": {
			Type: "number", Description: cdsUsage, Minimum: new(float64(0)),
		},
		"cue_insertion_offset_ms": {
			Type: "number", Description: ciomUsage, Minimum: new(float64(0)),
		},
		"cue_walltime_ms": {
			Type: "number", Description: cwmUsage, Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: pkg.OBOCOCUsage,
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
			Name: insertCuepointTool, Title: insertCuepointShort, Description: insertCuepointLong,
			InputSchema: insertCuepointInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertCuepointTool, cobramcp.ConfirmThen(
				func(input liveBroadcast.LiveBroadcast) string {
					return fmt.Sprintf(
						insertCuepointConfirm, strings.Join(input.Ids, ", "),
					)
				},
				func(input liveBroadcast.LiveBroadcast, w io.Writer) error {
					return input.InsertCuepoint(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(newInsertCuepointCmd())
}

type insertCuepointOpts struct {
	ids                  []string
	cueType              string
	cueDurationSecs      int64
	cueInsertionOffsetMs int64
	cueWalltimeMs        uint64
	oboco                string
	obococ               string
	output               string
	yes                  bool
}

func (o *insertCuepointOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, "IDs of the broadcasts")
	c.Flags().StringVarP(&o.cueType, "cueType", "c", "cueTypeAd", ctUsage)
	c.Flags().Int64VarP(
		&o.cueDurationSecs, "cueDurationSecs", "D", 0, cdsUsage,
	)
	c.Flags().Int64VarP(
		&o.cueInsertionOffsetMs, "cueInsertionOffsetMs", "O", 0, ciomUsage,
	)
	c.Flags().Uint64VarP(&o.cueWalltimeMs, "cueWalltimeMs", "W", 0, cwmUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *insertCuepointOpts) build() *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithIds(o.ids),
		liveBroadcast.WithCueType(o.cueType),
		liveBroadcast.WithCueDurationSecs(o.cueDurationSecs),
		liveBroadcast.WithCueInsertionOffsetMs(o.cueInsertionOffsetMs),
		liveBroadcast.WithCueWalltimeMs(o.cueWalltimeMs),
		liveBroadcast.WithOBOCO(o.oboco),
		liveBroadcast.WithOBOCOC(o.obococ),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
}

func newInsertCuepointCmd() *cobra.Command {
	var opts insertCuepointOpts

	c := &cobra.Command{
		Use:     "insertCuepoint",
		Short:   insertCuepointShort,
		Long:    insertCuepointLong,
		Example: insertCuepointExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertCuepointConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().InsertCuepoint(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

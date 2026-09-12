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
	insertCuepointConfirm = "Insert cuepoint into live broadcast(s): %s"
	insertCuepointShort   = "Insert a cuepoint into a live broadcast"
	insertCuepointLong    = "Insert a cuepoint into a live broadcast. Use this tool to insert an ad break cuepoint into a currently live broadcast."
	insertCuepointExample = `# Insert an ad cuepoint of 30 seconds
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
	liveBroadcastCmd.AddCommand(insertCuepointCmd)

	insertCuepointCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, "IDs of the broadcasts",
	)
	insertCuepointCmd.Flags().StringVarP(
		&cueType, "cueType", "c", "cueTypeAd", ctUsage,
	)
	insertCuepointCmd.Flags().Int64VarP(
		&cueDurationSecs, "cueDurationSecs", "D", 0, cdsUsage,
	)
	insertCuepointCmd.Flags().Int64VarP(
		&cueInsertionOffsetMs, "cueInsertionOffsetMs", "O", 0, ciomUsage,
	)
	insertCuepointCmd.Flags().Uint64VarP(
		&cueWalltimeMs, "cueWalltimeMs", "W", 0, cwmUsage,
	)
	insertCuepointCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCuepointCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	insertCuepointCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCuepointCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCuepointCmd.MarkFlagRequired("ids")
}

var insertCuepointCmd = &cobra.Command{
	Use:     "insertCuepoint",
	Short:   insertCuepointShort,
	Long:    insertCuepointLong,
	Example: insertCuepointExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(
			c, fmt.Sprintf(insertCuepointConfirm, strings.Join(ids, ", ")),
		)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithCueType(cueType),
			liveBroadcast.WithCueDurationSecs(cueDurationSecs),
			liveBroadcast.WithCueInsertionOffsetMs(cueInsertionOffsetMs),
			liveBroadcast.WithCueWalltimeMs(cueWalltimeMs),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.InsertCuepoint(c.OutOrStdout()), c)
	},
}

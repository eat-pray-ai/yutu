// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json"
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
		"cue_duration_secs":          {Type: "number", Description: cdsUsage},
		"cue_insertion_offset_ms":    {Type: "number", Description: ciomUsage},
		"cue_walltime_ms":            {Type: "number", Description: cwmUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
		}, cobramcp.GenToolHandler(
			insertCuepointTool,
			func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.InsertCuepoint(writer)
			},
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

	_ = insertCuepointCmd.MarkFlagRequired("ids")
	cmd.AddMutationFlags(insertCuepointCmd)
}

var insertCuepointCmd = &cobra.Command{
	Use:     "insertCuepoint",
	Short:   insertCuepointShort,
	Long:    insertCuepointLong,
	Example: insertCuepointExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would insert cuepoint into live broadcast(s): %s",
			strings.Join(ids, ", "),
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
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

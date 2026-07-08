// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"encoding/json"
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
)

const (
	updateTool    = "liveStream-update"
	updateIdUsage = "ID of the live stream to update"
	updateShort   = "Update a live stream"
	updateLong    = "Update a live stream. Use this tool to update an existing live stream's settings."
	updateExample = `# Update live stream title
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
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
		}, cobramcp.GenToolHandler(
			updateTool, func(input liveStream.LiveStream, writer io.Writer) error {
				input.MaxResults = 1
				return input.Update(writer)
			},
		),
	)
	liveStreamCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringVarP(&frameRate, "frameRate", "f", "", frUsage)
	updateCmd.Flags().StringVarP(&ingestionType, "ingestionType", "I", "", itUsage)
	updateCmd.Flags().StringVarP(&resolution, "resolution", "r", "", resUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	updateCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = updateCmd.MarkFlagRequired("id")
	cmd.AddMutationFlags(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would update live stream: %s", strings.Join(ids, ", "),
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveStream.NewLiveStream(
			liveStream.WithIds(ids),
			liveStream.WithTitle(title),
			liveStream.WithDescription(description),
			liveStream.WithFrameRate(frameRate),
			liveStream.WithIngestionType(ingestionType),
			liveStream.WithResolution(resolution),
			liveStream.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveStream.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveStream.WithMaxResults(1),
			liveStream.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}

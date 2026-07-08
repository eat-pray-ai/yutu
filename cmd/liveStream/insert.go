// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"encoding/json"
	"io"

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
	insertTool    = "liveStream-insert"
	insertShort   = "Insert a live stream"
	insertLong    = "Insert a live stream. Use this tool to create a new live stream for the authenticated user."
	insertExample = `# Create a live stream with RTMP at 1080p
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
			Default: json.RawMessage(`["snippet","cdn","status"]`),
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			insertTool, func(input liveStream.LiveStream, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	liveStreamCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&frameRate, "frameRate", "f", "", frUsage)
	insertCmd.Flags().StringVarP(&ingestionType, "ingestionType", "I", "", itUsage)
	insertCmd.Flags().StringVarP(&resolution, "resolution", "r", "", resUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "cdn", "status"}, pkg.PartsUsage,
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("title")
	_ = insertCmd.MarkFlagRequired("ingestionType")
	_ = insertCmd.MarkFlagRequired("resolution")
	cmd.AddMutationFlags(insertCmd)
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(c, "Would create live stream: %s", title)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveStream.NewLiveStream(
			liveStream.WithTitle(title),
			liveStream.WithDescription(description),
			liveStream.WithFrameRate(frameRate),
			liveStream.WithIngestionType(ingestionType),
			liveStream.WithResolution(resolution),
			liveStream.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveStream.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveStream.WithParts(parts),
			liveStream.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

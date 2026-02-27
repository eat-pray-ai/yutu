// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool  = "caption-update"
	updateShort = "Update caption"
	updateLong  = "Update caption of a video\n\nExamples:\n  yutu caption update --videoId dQw4w9WgXcQ --isDraft=false\n  yutu caption update --videoId dQw4w9WgXcQ --language en --name English\n  yutu caption update --videoId dQw4w9WgXcQ --file updated.srt"
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file": {Type: "string", Description: fileUsage},
		"audio_track_type": {
			Type: "string", Description: attUsage,
			Default: json.RawMessage(`"unknown"`),
		},
		"is_auto_synced": {Type: "boolean", Description: iasUsage},
		"is_cc":          {Type: "boolean", Description: iscUsage},
		"is_draft":       {Type: "boolean", Description: isdUsage},
		"is_easy_reader": {Type: "boolean", Description: iserUsage},
		"is_large":       {Type: "boolean", Description: islUsage},
		"language":       {Type: "string", Description: languageUsage},
		"name":           {Type: "string", Description: nameUsage},
		"track_kind": {
			Type: "string", Description: tkUsage,
			Enum:    []any{"standard", "ASR", "forced"},
			Default: json.RawMessage(`"standard"`),
		},
		"video_id":                   {Type: "string", Description: vidUsage},
		"on_behalf_of":               {Type: "string"},
		"on_behalf_of_content_owner": {Type: "string"},
		"output": {
			Type: "string", Description: pkg.SilentUsage,
			Enum:    []any{"json", "yaml", "silent", ""},
			Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
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
		}, cmd.GenToolHandler(
			updateTool, func(input caption.Caption, writer io.Writer) error {
				return input.Update(writer)
			},
		),
	)
	captionCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	updateCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	updateCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	updateCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	updateCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	updateCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	updateCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	updateCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	updateCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	updateCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("videoId")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := caption.NewCaption(
			caption.WithFile(file),
			caption.WithAudioTrackType(audioTrackType),
			caption.WithIsAutoSynced(isAutoSynced),
			caption.WithIsCC(isCC),
			caption.WithIsDraft(isDraft),
			caption.WithIsEasyReader(isEasyReader),
			caption.WithIsLarge(isLarge),
			caption.WithLanguage(language),
			caption.WithName(name),
			caption.WithTrackKind(trackKind),
			caption.WithVideoId(videoId),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithOutput(output),
			caption.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Update(cmd.OutOrStdout()), cmd)
	},
}

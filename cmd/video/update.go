// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool      = "video-update"
	updateShort     = "Update a video on YouTube"
	updateLong      = "Update a video on YouTube, with the specified title, description, tags, etc\n\nExamples:\n  yutu video update --id dQw4w9WgXcQ --title 'New Title'\n  yutu video update --id dQw4w9WgXcQ --description 'Updated description' --privacy public\n  yutu video update --id dQw4w9WgXcQ --tags 'music,pop,2024' --categoryId 10"
	updateIdUsage   = "ID of the video to update"
	updateLangUsage = "Language of the video"
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
		"tags": {
			Type: "array", Description: tagsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"language": {Type: "string", Description: updateLangUsage},
		"license": {
			Type: "string", Enum: []any{"youtube", "creativeCommon"},
			Description: licenseUsage, Default: json.RawMessage(`"youtube"`),
		},
		"thumbnail":   {Type: "string", Description: thumbnailUsage},
		"playlist_id": {Type: "string", Description: pidUsage},
		"category_id": {Type: "string", Description: caidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted", ""},
		},
		"embeddable": {Type: "boolean", Description: embeddableUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
			updateTool, func(input video.Video, writer io.Writer) error {
				return input.Update(writer)
			},
		),
	)
	videoCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", updateLangUsage)
	updateCmd.Flags().StringVarP(&license, "license", "L", "youtube", licenseUsage)
	updateCmd.Flags().StringVarP(&thumbnail, "thumbnail", "u", "", thumbnailUsage)
	updateCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", pidUsage)
	updateCmd.Flags().StringVarP(&categoryId, "categoryId", "g", "", caidUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().BoolVarP(
		embeddable, "embeddable", "E", true, embeddableUsage,
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithTitle(title),
			video.WithDescription(description),
			video.WithTags(tags),
			video.WithLanguage(language),
			video.WithLicense(license),
			video.WithPlaylistId(playListId),
			video.WithThumbnail(thumbnail),
			video.WithCategory(categoryId),
			video.WithPrivacy(privacy),
			video.WithEmbeddable(embeddable),
			video.WithOutput(output),
			video.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Update(cmd.OutOrStdout()), cmd)
	},
}

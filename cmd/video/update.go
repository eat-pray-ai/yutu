// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

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
	updateLong      = "Update a video on YouTube, with the specified title, description, tags, etc"
	updateIdUsage   = "ID of the video to update"
	updateLangUsage = "Language of the video"
)

type updateIn struct {
	Ids         []string `json:"ids"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Language    string   `json:"language"`
	License     string   `json:"license"`
	Thumbnail   string   `json:"thumbnail"`
	PlaylistId  string   `json:"playlistId"`
	CategoryId  string   `json:"categoryId"`
	Privacy     string   `json:"privacy"`
	Embeddable  *string  `json:"embeddable,omitempty"`
	Output      string   `json:"output"`
	Jsonpath    string   `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
			Default:     json.RawMessage(`[]`),
		},
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"tags": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: tagsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"language": {
			Type: "string", Description: updateLangUsage,
			Default: json.RawMessage(`""`),
		},
		"license": {
			Type: "string", Enum: []any{"youtube", "creativeCommon"},
			Description: licenseUsage, Default: json.RawMessage(`"youtube"`),
		},
		"thumbnail": {
			Type: "string", Description: thumbnailUsage,
			Default: json.RawMessage(`""`),
		},
		"playlistId": {
			Type: "string", Description: pidUsage,
			Default: json.RawMessage(`""`),
		},
		"categoryId": {
			Type: "string", Description: caidUsage,
			Default: json.RawMessage(`""`),
		},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage, Default: json.RawMessage(`""`),
		},
		"embeddable": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: embeddableUsage, Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
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
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, req *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: updateTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.Ids
	title = input.Title
	description = input.Description
	tags = input.Tags
	language = input.Language
	license = input.License
	thumbnail = input.Thumbnail
	playListId = input.PlaylistId
	categoryId = input.CategoryId
	privacy = input.Privacy
	embeddable = utils.StrToBoolPtr(input.Embeddable)
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func update(writer io.Writer) error {
	v := video.NewVideo(
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
		video.WithMaxResults(1),
		video.WithService(nil),
	)

	return v.Update(output, jsonpath, writer)
}

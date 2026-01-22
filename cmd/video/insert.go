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
	insertTool      = "video-insert"
	insertShort     = "Upload a video to YouTube"
	insertLong      = "Upload a video to YouTube, with the specified title, description, tags, etc"
	insertLangUsage = "Language of the video"
)

type insertIn struct {
	AutoLevels                    *string  `json:"autoLevels,omitempty"`
	File                          string   `json:"file"`
	Title                         string   `json:"title"`
	Description                   string   `json:"description"`
	Tags                          []string `json:"tags"`
	Language                      string   `json:"language"`
	License                       string   `json:"license"`
	Thumbnail                     string   `json:"thumbnail"`
	ChannelId                     string   `json:"channelId"`
	PlaylistId                    string   `json:"playlistId"`
	CategoryId                    string   `json:"categoryId"`
	Privacy                       string   `json:"privacy"`
	ForKids                       *string  `json:"forKids,omitempty"`
	Embeddable                    *string  `json:"embeddable,omitempty"`
	PublishAt                     string   `json:"publishAt"`
	Stabilize                     *string  `json:"stabilize,omitempty"`
	NotifySubscribers             *string  `json:"notifySubscribers,omitempty"`
	PublicStatsViewable           *string  `json:"publicStatsViewable,omitempty"`
	OnBehalfOfContentOwner        string   `json:"onBehalfOfContentOwner"`
	OnBehalfOfContentOwnerChannel string   `json:"onBehalfOfContentOwnerChannel"`
	Output                        string   `json:"output"`
	Jsonpath                      string   `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "categoryId", "privacy"},
	Properties: map[string]*jsonschema.Schema{
		"autoLevels": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: alUsage, Default: json.RawMessage(`""`),
		},
		"file": {
			Type: "string", Description: fileUsage,
			Default: json.RawMessage(`""`),
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
			Type: "string", Description: insertLangUsage,
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
		"channelId": {
			Type: "string", Description: chidUsage,
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
		"forKids": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: fkUsage, Default: json.RawMessage(`""`),
		},
		"embeddable": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: embeddableUsage, Default: json.RawMessage(`""`),
		},
		"publishAt": {
			Type: "string", Description: paUsage,
			Default: json.RawMessage(`""`),
		},
		"stabilize": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: stabilizeUsage, Default: json.RawMessage(`""`),
		},
		"notifySubscribers": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: nsUsage, Default: json.RawMessage(`""`),
		},
		"publicStatsViewable": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: psvUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwnerChannel": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().BoolVarP(
		autoLevels, "autoLevels", "A", true, alUsage,
	)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", insertLangUsage)
	insertCmd.Flags().StringVarP(&license, "license", "L", "youtube", licenseUsage)
	insertCmd.Flags().StringVarP(&thumbnail, "thumbnail", "u", "", thumbnailUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", chidUsage)
	insertCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", pidUsage)
	insertCmd.Flags().StringVarP(&categoryId, "categoryId", "g", "", caidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().BoolVarP(forKids, "forKids", "K", false, fkUsage)
	insertCmd.Flags().BoolVarP(
		embeddable, "embeddable", "E", true, embeddableUsage,
	)
	insertCmd.Flags().StringVarP(&publishAt, "publishAt", "U", "", paUsage)
	insertCmd.Flags().BoolVarP(stabilize, "stabilize", "S", true, stabilizeUsage)
	insertCmd.Flags().BoolVarP(
		notifySubscribers, "notifySubscribers", "N", true, nsUsage,
	)
	insertCmd.Flags().BoolVarP(
		publicStatsViewable, "publicStatsViewable", "P", false, psvUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("categoryId")
	_ = insertCmd.MarkFlagRequired("privacy")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	autoLevels = utils.BoolPtr(*input.AutoLevels)
	file = input.File
	title = input.Title
	description = input.Description
	tags = input.Tags
	language = input.Language
	license = input.License
	thumbnail = input.Thumbnail
	channelId = input.ChannelId
	playListId = input.PlaylistId
	categoryId = input.CategoryId
	privacy = input.Privacy
	forKids = utils.BoolPtr(*input.ForKids)
	embeddable = utils.BoolPtr(*input.Embeddable)
	publishAt = input.PublishAt
	stabilize = utils.BoolPtr(*input.Stabilize)
	notifySubscribers = utils.BoolPtr(*input.NotifySubscribers)
	publicStatsViewable = utils.BoolPtr(*input.PublicStatsViewable)
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	onBehalfOfContentOwnerChannel = input.OnBehalfOfContentOwnerChannel
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	v := video.NewVideo(
		video.WithAutoLevels(autoLevels),
		video.WithFile(file),
		video.WithTitle(title),
		video.WithDescription(description),
		video.WithTags(tags),
		video.WithLanguage(language),
		video.WithLicense(license),
		video.WithThumbnail(thumbnail),
		video.WithChannelId(channelId),
		video.WithPlaylistId(playListId),
		video.WithCategory(categoryId),
		video.WithPrivacy(privacy),
		video.WithForKids(forKids),
		video.WithEmbeddable(embeddable),
		video.WithPublishAt(publishAt),
		video.WithStabilize(stabilize),
		video.WithNotifySubscribers(notifySubscribers),
		video.WithPublicStatsViewable(publicStatsViewable),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		video.WithService(nil),
	)

	return v.Insert(output, jsonpath, writer)
}

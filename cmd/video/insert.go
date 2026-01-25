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
			Type: "string", Enum: []any{"true", "false", ""}, Description: alUsage,
		},
		"file":        {Type: "string", Description: fileUsage},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"tags": {
			Type: "array", Description: tagsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"language": {Type: "string", Description: insertLangUsage},
		"license": {
			Type: "string", Enum: []any{"youtube", "creativeCommon"},
			Description: licenseUsage, Default: json.RawMessage(`"youtube"`),
		},
		"thumbnail":  {Type: "string", Description: thumbnailUsage},
		"channelId":  {Type: "string", Description: chidUsage},
		"playlistId": {Type: "string", Description: pidUsage},
		"categoryId": {Type: "string", Description: caidUsage},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage,
		},
		"forKids": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: fkUsage,
		},
		"embeddable": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: embeddableUsage,
		},
		"publishAt": {
			Type: "string", Description: paUsage,
		},
		"stabilize": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: stabilizeUsage,
		},
		"notifySubscribers": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: nsUsage,
		},
		"publicStatsViewable": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: psvUsage,
		},
		"onBehalfOfContentOwner":        {Type: "string"},
		"onBehalfOfContentOwnerChannel": {Type: "string"},
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

	autoLevels = utils.StrToBoolPtr(input.AutoLevels)
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
	forKids = utils.StrToBoolPtr(input.ForKids)
	embeddable = utils.StrToBoolPtr(input.Embeddable)
	publishAt = input.PublishAt
	stabilize = utils.StrToBoolPtr(input.Stabilize)
	notifySubscribers = utils.StrToBoolPtr(input.NotifySubscribers)
	publicStatsViewable = utils.StrToBoolPtr(input.PublicStatsViewable)
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

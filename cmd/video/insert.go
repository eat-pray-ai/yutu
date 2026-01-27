// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
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

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "category_id", "privacy"},
	Properties: map[string]*jsonschema.Schema{
		"auto_levels": {Type: "boolean", Description: alUsage},
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
		"thumbnail":   {Type: "string", Description: thumbnailUsage},
		"channel_id":  {Type: "string", Description: chidUsage},
		"playlist_id": {Type: "string", Description: pidUsage},
		"category_id": {Type: "string", Description: caidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted", ""},
		},
		"for_kids":              {Type: "boolean", Description: fkUsage},
		"embeddable":            {Type: "boolean", Description: embeddableUsage},
		"publish_at":            {Type: "string", Description: paUsage},
		"stabilize":             {Type: "boolean", Description: stabilizeUsage},
		"notify_subscribers":    {Type: "boolean", Description: nsUsage},
		"public_stats_viewable": {Type: "boolean", Description: psvUsage},

		"on_behalf_of_content_owner":         {Type: "string"},
		"on_behalf_of_content_owner_channel": {Type: "string"},
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
		}, cmd.GenToolHandler(
			insertTool, func(input video.Video, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
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
		input := video.NewVideo(
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
			video.WithOutput(output),
			video.WithJsonpath(jsonpath),
		)
		err := input.Insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

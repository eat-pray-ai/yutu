// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	updateTool      = "video-update"
	updateIdUsage   = "ID of the video to update"
	updateLangUsage = "Language of the video"
	updateConfirm   = "Update video: %s"
	updateShort     = "Update a video"
	updateLong      = "Update a video. Use this tool to update a video."
	updateExample   = `  # Update video title
  yutu video update --id dQw4w9WgXcQ --title 'New Title'
  # Update video description and privacy
  yutu video update --id dQw4w9WgXcQ --description 'Updated description' --privacy public
  # Schedule a private video to publish later
  yutu video update --id dQw4w9WgXcQ --privacy private --publishAt 2026-08-18T15:00:00Z
  # Update video tags and category
  yutu video update --id dQw4w9WgXcQ --tags 'music,pop,2024' --categoryId 10`
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
			Description: licenseUsage,
		},
		"thumbnail":   {Type: "string", Description: thumbnailUsage},
		"playlist_id": {Type: "string", Description: pidUsage},
		"category_id": {Type: "string", Description: caidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
		"publish_at":               {Type: "string", Description: paUsage},
		"embeddable":               {Type: "boolean", Description: embeddableUsage},
		"contains_synthetic_media": {Type: "boolean", Description: csmUsage},
		"recording_date":           {Type: "string", Description: rdUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
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
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input video.Video) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input video.Video, w io.Writer) error {
					video.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	videoCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids                    []string
	title                  string
	description            string
	tags                   []string
	language               string
	license                string
	thumbnail              string
	playlistId             string
	categoryId             string
	privacy                string
	publishAt              string
	embeddable             bool
	containsSyntheticMedia bool
	recordingDate          string
	output                 string
	yes                    bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringSliceVarP(&o.tags, "tags", "a", nil, tagsUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", updateLangUsage)
	c.Flags().StringVarP(&o.license, "license", "L", "", licenseUsage)
	c.Flags().StringVarP(&o.thumbnail, "thumbnail", "u", "", thumbnailUsage)
	c.Flags().StringVarP(&o.playlistId, "playlistId", "y", "", pidUsage)
	c.Flags().StringVarP(&o.categoryId, "categoryId", "g", "", caidUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().StringVarP(&o.publishAt, "publishAt", "U", "", paUsage)
	c.Flags().BoolVarP(&o.embeddable, "embeddable", "E", true, embeddableUsage)
	c.Flags().BoolVarP(
		&o.containsSyntheticMedia, "containsSyntheticMedia", "M", false, csmUsage,
	)
	c.Flags().StringVarP(&o.recordingDate, "recordingDate", "D", "", rdUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *video.Video {
	return video.NewVideo(
		video.WithIds(o.ids),
		video.WithTitle(o.title),
		utils.OptIfChanged(f, "description", &o.description, video.WithDescription),
		utils.OptIfChanged(f, "tags", o.tags, video.WithTags),
		utils.OptIfChanged(f, "language", &o.language, video.WithLanguage),
		video.WithLicense(o.license),
		video.WithPlaylistId(o.playlistId),
		video.WithThumbnail(o.thumbnail),
		video.WithCategory(o.categoryId),
		video.WithPrivacy(o.privacy),
		utils.OptIfChanged(f, "publishAt", &o.publishAt, video.WithPublishAt),
		utils.OptIfChanged(f, "embeddable", &o.embeddable, video.WithEmbeddable),
		utils.OptIfChanged(
			f, "containsSyntheticMedia", &o.containsSyntheticMedia,
			video.WithContainsSyntheticMedia,
		),
		utils.OptIfChanged(
			f, "recordingDate", &o.recordingDate, video.WithRecordingDate,
		),
		video.WithMaxResults(1),
		video.WithOutput(o.output),
	).(*video.Video)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

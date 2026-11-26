// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

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
	insertTool      = "video-insert"
	insertLangUsage = "Language of the video"
	insertConfirm   = "Upload video: %s"
	insertShort     = "Upload a video"
	insertLong      = "Upload a video. Use this tool to upload a video."
	insertExample   = `  # Upload a public video
  yutu video insert --file video.mp4 --title 'My Video' --categoryId 22 --privacy public
  # Upload a private video with tags
  yutu video insert --file video.mp4 --title 'Tutorial' --categoryId 27 --privacy private --tags 'go,tutorial'
  # Upload an unlisted video with custom thumbnail
  yutu video insert --file video.mp4 --title 'Music Video' --categoryId 10 --privacy unlisted --thumbnail cover.jpg`
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
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["yutu🐰"]`),
		},
		"language": {Type: "string", Description: insertLangUsage},
		"license": {
			Type: "string", Enum: []any{"youtube", "creativeCommon"},
			Description: licenseUsage,
		},
		"thumbnail":   {Type: "string", Description: thumbnailUsage},
		"channel_id":  {Type: "string", Description: chidUsage},
		"playlist_id": {Type: "string", Description: pidUsage},
		"category_id": {Type: "string", Description: caidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
		"for_kids": {
			Type: "boolean", Description: fkUsage, Default: jsontext.Value("false"),
		},
		"contains_synthetic_media": {
			Type: "boolean", Description: csmUsage, Default: jsontext.Value("false"),
		},
		"embeddable":            {Type: "boolean", Description: embeddableUsage},
		"recording_date":        {Type: "string", Description: rdUsage},
		"publish_at":            {Type: "string", Description: paUsage},
		"stabilize":             {Type: "boolean", Description: stabilizeUsage},
		"notify_subscribers":    {Type: "boolean", Description: nsUsage},
		"public_stats_viewable": {Type: "boolean", Description: psvUsage},
		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
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
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input video.Video) string {
					return fmt.Sprintf(insertConfirm, input.File)
				},
				func(input video.Video, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	videoCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	autoLevels                    bool
	file                          string
	title                         string
	description                   string
	tags                          []string
	language                      string
	license                       string
	thumbnail                     string
	channelId                     string
	playlistId                    string
	categoryId                    string
	privacy                       string
	forKids                       bool
	embeddable                    bool
	containsSyntheticMedia        bool
	recordingDate                 string
	publishAt                     string
	stabilize                     bool
	notifySubscribers             bool
	publicStatsViewable           bool
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	output                        string
	yes                           bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().BoolVarP(&o.autoLevels, "autoLevels", "A", true, alUsage)
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringSliceVarP(&o.tags, "tags", "a", []string{"yutu🐰"}, tagsUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", insertLangUsage)
	c.Flags().StringVarP(&o.license, "license", "L", "", licenseUsage)
	c.Flags().StringVarP(&o.thumbnail, "thumbnail", "u", "", thumbnailUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", chidUsage)
	c.Flags().StringVarP(&o.playlistId, "playlistId", "y", "", pidUsage)
	c.Flags().StringVarP(&o.categoryId, "categoryId", "g", "", caidUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().BoolVarP(&o.forKids, "forKids", "K", false, fkUsage)
	c.Flags().BoolVarP(&o.embeddable, "embeddable", "E", true, embeddableUsage)
	c.Flags().BoolVarP(
		&o.containsSyntheticMedia, "containsSyntheticMedia", "M", false, csmUsage,
	)
	c.Flags().StringVarP(&o.recordingDate, "recordingDate", "D", "", rdUsage)
	c.Flags().StringVarP(&o.publishAt, "publishAt", "U", "", paUsage)
	c.Flags().BoolVarP(&o.stabilize, "stabilize", "S", true, stabilizeUsage)
	c.Flags().BoolVarP(&o.notifySubscribers, "notifySubscribers", "N", true, nsUsage)
	c.Flags().BoolVarP(
		&o.publicStatsViewable, "publicStatsViewable", "P", false, psvUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("file")
	_ = c.MarkFlagRequired("categoryId")
	_ = c.MarkFlagRequired("privacy")
}

func (o *insertOpts) build(f *pflag.FlagSet) *video.Video {
	return video.NewVideo(
		utils.OptIfChanged(f, "autoLevels", &o.autoLevels, video.WithAutoLevels),
		video.WithFile(o.file),
		video.WithTitle(o.title),
		utils.OptIfChanged(f, "description", &o.description, video.WithDescription),
		video.WithTags(o.tags),
		utils.OptIfChanged(f, "language", &o.language, video.WithLanguage),
		video.WithLicense(o.license),
		video.WithThumbnail(o.thumbnail),
		video.WithChannelId(o.channelId),
		video.WithPlaylistId(o.playlistId),
		video.WithCategory(o.categoryId),
		video.WithPrivacy(o.privacy),
		utils.OptIfChanged(f, "forKids", &o.forKids, video.WithForKids),
		utils.OptIfChanged(f, "embeddable", &o.embeddable, video.WithEmbeddable),
		utils.OptIfChanged(
			f, "containsSyntheticMedia", &o.containsSyntheticMedia,
			video.WithContainsSyntheticMedia,
		),
		utils.OptIfChanged(
			f, "recordingDate", &o.recordingDate, video.WithRecordingDate,
		),
		utils.OptIfChanged(f, "publishAt", &o.publishAt, video.WithPublishAt),
		utils.OptIfChanged(f, "stabilize", &o.stabilize, video.WithStabilize),
		utils.OptIfChanged(
			f, "notifySubscribers", &o.notifySubscribers,
			video.WithNotifySubscribers,
		),
		utils.OptIfChanged(
			f, "publicStatsViewable", &o.publicStatsViewable,
			video.WithPublicStatsViewable,
		),
		video.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		video.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		video.WithOutput(o.output),
	).(*video.Video)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.file))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

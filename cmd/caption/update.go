// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	updateTool    = "caption-update"
	updateConfirm = "Update caption for video %s"
	updateShort   = "Update a caption"
	updateLong    = "Update a caption. Use this tool to update a caption."
	updateExample = `  # Publish a draft caption
  yutu caption update --videoId dQw4w9WgXcQ --isDraft=false
  # Update caption language and name
  yutu caption update --videoId dQw4w9WgXcQ --language en --name English
  # Update caption file
  yutu caption update --videoId dQw4w9WgXcQ --file updated.srt`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file": {Type: "string", Description: fileUsage},
		"audio_track_type": {
			Type: "string", Description: attUsage,
			Enum: []any{"unknown", "primary", "commentary", "descriptive"},
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
			Enum: []any{"standard", "ASR", "forced"},
		},
		"video_id":                   {Type: "string", Description: vidUsage},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
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
				func(input caption.Caption) string {
					target := input.VideoId
					if input.File != "" {
						target += " with file: " + input.File
					}
					return fmt.Sprintf(updateConfirm, target)
				},
				func(input caption.Caption, w io.Writer) error {
					return input.Update(w)
				},
			),
		),
	)
	captionCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	file                   string
	audioTrackType         string
	isAutoSynced           bool
	isCC                   bool
	isDraft                bool
	isEasyReader           bool
	isLarge                bool
	language               string
	name                   string
	trackKind              string
	videoId                string
	onBehalfOf             string
	onBehalfOfContentOwner string
	output                 string
	yes                    bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(&o.audioTrackType, "audioTrackType", "a", "", attUsage)
	c.Flags().BoolVarP(&o.isAutoSynced, "isAutoSynced", "A", true, iasUsage)
	c.Flags().BoolVarP(&o.isCC, "isCC", "C", false, iscUsage)
	c.Flags().BoolVarP(&o.isDraft, "isDraft", "D", false, isdUsage)
	c.Flags().BoolVarP(&o.isEasyReader, "isEasyReader", "E", false, iserUsage)
	c.Flags().BoolVarP(&o.isLarge, "isLarge", "L", false, islUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", languageUsage)
	c.Flags().StringVarP(&o.name, "name", "n", "", nameUsage)
	c.Flags().StringVarP(&o.trackKind, "trackKind", "t", "", tkUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(&o.onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("videoId")
}

func (o *updateOpts) build(f *pflag.FlagSet) *caption.Caption {
	return caption.NewCaption(
		caption.WithFile(o.file),
		utils.OptIfChanged(
			f, "audioTrackType", o.audioTrackType, caption.WithAudioTrackType,
		),
		utils.OptIfChanged(f, "isAutoSynced", &o.isAutoSynced, caption.WithIsAutoSynced),
		utils.OptIfChanged(f, "isCC", &o.isCC, caption.WithIsCC),
		utils.OptIfChanged(f, "isDraft", &o.isDraft, caption.WithIsDraft),
		utils.OptIfChanged(f, "isEasyReader", &o.isEasyReader, caption.WithIsEasyReader),
		utils.OptIfChanged(f, "isLarge", &o.isLarge, caption.WithIsLarge),
		caption.WithLanguage(o.language),
		caption.WithName(o.name),
		utils.OptIfChanged(f, "trackKind", o.trackKind, caption.WithTrackKind),
		caption.WithVideoId(o.videoId),
		caption.WithOnBehalfOf(o.onBehalfOf),
		caption.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		caption.WithOutput(o.output),
	).(*caption.Caption)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			target := opts.videoId
			if opts.file != "" {
				target += " with file: " + opts.file
			}
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(updateConfirm, target))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

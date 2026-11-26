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
	insertTool    = "caption-insert"
	insertConfirm = "Upload caption for video %s with %s"
	insertShort   = "Upload a caption"
	insertLong    = "Upload a caption. Use this tool to upload a caption to a video."
	insertExample = `  # Insert a caption to a video
  yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ
  # Insert a caption with language and name
  yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ --language en --name English
  # Insert a standard caption and publish immediately
  yutu caption insert --file subtitle.srt --videoId dQw4w9WgXcQ --isDraft=false`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":                       {Type: "string", Description: fileUsage},
		"is_auto_synced":             {Type: "boolean", Description: iasUsage},
		"is_cc":                      {Type: "boolean", Description: iscUsage},
		"is_draft":                   {Type: "boolean", Description: isdUsage},
		"is_easy_reader":             {Type: "boolean", Description: iserUsage},
		"is_large":                   {Type: "boolean", Description: islUsage},
		"language":                   {Type: "string", Description: languageUsage},
		"name":                       {Type: "string", Description: nameUsage},
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input caption.Caption) string {
					return fmt.Sprintf(insertConfirm, input.VideoId, input.File)
				},
				func(input caption.Caption, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	captionCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	file         string
	isAutoSynced bool
	isCC         bool
	isDraft      bool
	isEasyReader bool
	isLarge      bool
	language     string
	name         string
	videoId      string
	obo          string
	oboco        string
	output       string
	yes          bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().BoolVarP(&o.isAutoSynced, "isAutoSynced", "A", true, iasUsage)
	c.Flags().BoolVarP(&o.isCC, "isCC", "C", false, iscUsage)
	c.Flags().BoolVarP(&o.isDraft, "isDraft", "D", false, isdUsage)
	c.Flags().BoolVarP(&o.isEasyReader, "isEasyReader", "E", false, iserUsage)
	c.Flags().BoolVarP(&o.isLarge, "isLarge", "L", false, islUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", languageUsage)
	c.Flags().StringVarP(&o.name, "name", "n", "", nameUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(&o.obo, "onBehalfOf", "B", "", pkg.OBOUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("file")
	_ = c.MarkFlagRequired("videoId")
}

func (o *insertOpts) build(f *pflag.FlagSet) *caption.Caption {
	return caption.NewCaption(
		caption.WithFile(o.file),
		utils.OptIfChanged(f, "isAutoSynced", &o.isAutoSynced, caption.WithIsAutoSynced),
		utils.OptIfChanged(f, "isCC", &o.isCC, caption.WithIsCC),
		utils.OptIfChanged(f, "isDraft", &o.isDraft, caption.WithIsDraft),
		utils.OptIfChanged(f, "isEasyReader", &o.isEasyReader, caption.WithIsEasyReader),
		utils.OptIfChanged(f, "isLarge", &o.isLarge, caption.WithIsLarge),
		caption.WithLanguage(o.language),
		caption.WithName(o.name),
		caption.WithVideoId(o.videoId),
		caption.WithOBO(o.obo),
		caption.WithOBOCO(o.oboco),
		caption.WithOutput(o.output),
	).(*caption.Caption)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.videoId, opts.file),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

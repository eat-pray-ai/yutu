// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	downloadTool    = "caption-download"
	downloadIdUsage = "ID of the caption to download"
	downloadShort   = "Download a caption"
	downloadLong    = "Download a caption. Use this tool to download a caption from a video."
	downloadExample = `  # Download a caption as SRT
  yutu caption download --id abc123 --file subtitle.srt
  # Download a caption as VTT
  yutu caption download --id abc123 --file subtitle.vtt --tfmt vtt
  # Download a caption translated to French
  yutu caption download --id abc123 --file subtitle.srt --tlang fr`
)

var downloadInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "file"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: downloadIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"file": {Type: "string", Description: fileUsage},
		"tfmt": {
			Type: "string", Description: tfmtUsage,
			Enum: []any{"sbv", "srt", "vtt"},
		},
		"tlang":                      {Type: "string", Description: tlangUsage},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: downloadTool, Title: downloadShort, Description: downloadLong,
			InputSchema: downloadInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			downloadTool, func(input caption.Caption, writer io.Writer) error {
				return input.Download(writer)
			},
		),
	)
	captionCmd.AddCommand(newDownloadCmd())
}

type downloadOpts struct {
	ids                    []string
	file                   string
	tfmt                   string
	tlang                  string
	onBehalfOf             string
	onBehalfOfContentOwner string
}

func (o *downloadOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, downloadIdUsage)
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(&o.tfmt, "tfmt", "t", "", tfmtUsage)
	c.Flags().StringVarP(&o.tlang, "tlang", "l", "", tlangUsage)
	c.Flags().StringVarP(&o.onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	_ = c.MarkFlagRequired("id")
	_ = c.MarkFlagRequired("file")
}

func (o *downloadOpts) build() *caption.Caption {
	return caption.NewCaption(
		caption.WithIds(o.ids),
		caption.WithFile(o.file),
		caption.WithTfmt(o.tfmt),
		caption.WithTlang(o.tlang),
		caption.WithOnBehalfOf(o.onBehalfOf),
		caption.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
	).(*caption.Caption)
}

func newDownloadCmd() *cobra.Command {
	var opts downloadOpts

	c := &cobra.Command{
		Use:     "download",
		Short:   downloadShort,
		Long:    downloadLong,
		Example: downloadExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Download(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

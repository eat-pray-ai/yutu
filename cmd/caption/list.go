// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"encoding/json/jsontext"
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
	listTool     = "caption-list"
	listIdsUsage = "IDs of the captions to list"
	listShort    = "List captions"
	listLong     = "List captions. Use this tool to list captions of a video."
	listExample  = `  # List captions of a video
  yutu caption list --videoId dQw4w9WgXcQ
  # List captions in JSON format
  yutu caption list --videoId dQw4w9WgXcQ --output json
  # List specific captions by IDs
  yutu caption list --ids abc123,def456 --videoId dQw4w9WgXcQ`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"video_id":                   {Type: "string", Description: vidUsage},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cobramcp.GenToolHandler(
			listTool, func(input caption.Caption, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	captionCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                    []string
	videoId                string
	onBehalfOf             string
	onBehalfOfContentOwner string
	parts                  []string
	output                 string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(&o.onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *caption.Caption {
	return caption.NewCaption(
		caption.WithIds(o.ids),
		caption.WithVideoId(o.videoId),
		caption.WithOnBehalfOf(o.onBehalfOf),
		caption.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		caption.WithParts(o.parts),
		caption.WithOutput(o.output),
	).(*caption.Caption)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

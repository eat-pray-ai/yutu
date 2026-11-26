// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "playlistImage-list"
	listShort   = "List playlist images"
	listLong    = "List playlist images. Use this tool to list playlist images."
	listExample = `  # List images of a playlist
  yutu playlistImage list --parent PLxxx
  # List playlist images with limit in JSON format
  yutu playlistImage list --parent PLxxx --maxResults 10 --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"parent": {Type: "string", Description: parentUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","kind","snippet"]`),
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
			listTool, func(input playlistImage.PlaylistImage, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	playlistImageCmd.AddCommand(newListCmd())
}

type listOpts struct {
	parent                        string
	maxResults                    int64
	parts                         []string
	output                        string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.parent, "parent", "P", "", parentUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "kind", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
}

func (o *listOpts) build() *playlistImage.PlaylistImage {
	return playlistImage.NewPlaylistImage(
		playlistImage.WithParent(o.parent),
		playlistImage.WithMaxResults(o.maxResults),
		playlistImage.WithParts(o.parts),
		playlistImage.WithOutput(o.output),
		playlistImage.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		playlistImage.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		playlistImage.WithService(nil),
	).(*playlistImage.PlaylistImage)
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

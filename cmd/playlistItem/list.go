// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "playlistItem-list"
	listIdsUsage = "IDs of the playlist items to list"
	listPidUsage = "Return the playlist items within the given playlist"
	listShort    = "List playlist items"
	listLong     = "List playlist items. Use this tool to list playlist items."
	listExample  = `  # List items in a playlist
  yutu playlistItem list --playlistId PLxxx
  # List playlist items with limit in JSON format
  yutu playlistItem list --playlistId PLxxx --maxResults 20 --output json
  # List specific playlist items by IDs
  yutu playlistItem list --ids abc123,def456`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"playlist_id": {Type: "string", Description: listPidUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"video_id":                   {Type: "string", Description: vidUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","status"]`),
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
			listTool, func(input playlistItem.PlaylistItem, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	playlistItemCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                    []string
	playlistId             string
	maxResults             int64
	videoId                string
	onBehalfOfContentOwner string
	parts                  []string
	output                 string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.playlistId, "playlistId", "y", "", listPidUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", vidUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *playlistItem.PlaylistItem {
	return playlistItem.NewPlaylistItem(
		playlistItem.WithIds(o.ids),
		playlistItem.WithPlaylistId(o.playlistId),
		playlistItem.WithMaxResults(o.maxResults),
		playlistItem.WithVideoId(o.videoId),
		playlistItem.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		playlistItem.WithParts(o.parts),
		playlistItem.WithOutput(o.output),
	).(*playlistItem.PlaylistItem)
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

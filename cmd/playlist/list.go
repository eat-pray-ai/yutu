// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	listTool     = "playlist-list"
	listIdsUsage = "Return the playlists with the given Ids for Stubby or Apiary"
	listCidUsage = "Return the playlists owned by the specified channel id"
	listShort    = "List playlists"
	listLong     = "List playlists. Use this tool to list playlists."
	listExample  = `  # List my playlists
  yutu playlist list --mine
  # List playlists by channel ID with limit
  yutu playlist list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10
  # List specific playlists in JSON format
  yutu playlist list --ids PLxxx1,PLxxx2 --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"channel_id": {Type: "string", Description: listCidUsage},
		"hl":         {Type: "string", Description: hlUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"mine": {
			Type: "boolean", Description: mineUsage,
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
			listTool, func(input playlist.Playlist, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	playlistCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                           []string
	channelId                     string
	hl                            string
	maxResults                    int64
	mine                          bool
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	parts                         []string
	output                        string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", listCidUsage)
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().BoolVarP(&o.mine, "mine", "M", true, mineUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build(f *pflag.FlagSet) *playlist.Playlist {
	return playlist.NewPlaylist(
		playlist.WithIds(o.ids),
		playlist.WithChannelId(o.channelId),
		playlist.WithHl(o.hl),
		playlist.WithMaxResults(o.maxResults),
		utils.OptIfChanged(f, "mine", &o.mine, playlist.WithMine),
		playlist.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		playlist.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		playlist.WithParts(o.parts),
		playlist.WithOutput(o.output),
	).(*playlist.Playlist)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"encoding/json/jsontext"
	"fmt"
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
	insertTool     = "playlist-insert"
	insertCidUsage = "Channel id of the playlist"
	insertConfirm  = "Create playlist: %s"
	insertShort    = "Create a playlist"
	insertLong     = "Create a playlist. Use this tool to create a playlist."
	insertExample  = `  # Create a public playlist
  yutu playlist insert --title 'My Playlist' --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --privacy public
  # Create a private playlist with description
  yutu playlist insert --title 'Tutorial Series' --channelId UC_x5X --privacy private --description 'My tutorials'
  # Create an unlisted playlist with tags
  yutu playlist insert --title 'Music' --channelId UC_x5X --privacy unlisted --tags 'music,pop'`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"title", "channel_id", "privacy"},
	Properties: map[string]*jsonschema.Schema{
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"tags": {
			Type: "array", Description: tagsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"language":   {Type: "string", Description: languageUsage},
		"channel_id": {Type: "string", Description: insertCidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
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
				func(input playlist.Playlist) string {
					return fmt.Sprintf(insertConfirm, input.Title)
				},
				func(input playlist.Playlist, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	playlistCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	title       string
	description string
	tags        []string
	language    string
	channelId   string
	privacy     string
	output      string
	yes         bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringSliceVarP(&o.tags, "tags", "a", nil, tagsUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", languageUsage)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", insertCidUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("channelId")
	_ = c.MarkFlagRequired("privacy")
}

func (o *insertOpts) build(f *pflag.FlagSet) *playlist.Playlist {
	return playlist.NewPlaylist(
		playlist.WithTitle(o.title),
		utils.OptIfChanged(f, "description", &o.description, playlist.WithDescription),
		utils.OptIfChanged(f, "tags", o.tags, playlist.WithTags),
		utils.OptIfChanged(f, "language", &o.language, playlist.WithLanguage),
		playlist.WithChannelId(o.channelId),
		playlist.WithPrivacy(o.privacy),
		playlist.WithOutput(o.output),
	).(*playlist.Playlist)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.title))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

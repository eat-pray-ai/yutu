// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

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
	updateTool    = "playlist-update"
	updateIdUsage = "ID of the playlist to update"
	updateConfirm = "Update playlist: %s"
	updateShort   = "Update a playlist"
	updateLong    = "Update a playlist. Use this tool to update a playlist."
	updateExample = `  # Update playlist title
  yutu playlist update --id PLxxx --title 'Updated Title'
  # Update playlist description and privacy
  yutu playlist update --id PLxxx --description 'New description' --privacy public
  # Update playlist tags
  yutu playlist update --id PLxxx --tags 'music,pop,2024'`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
		},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"tags": {
			Type: "array", Description: tagsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"language": {Type: "string", Description: languageUsage},
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input playlist.Playlist) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input playlist.Playlist, w io.Writer) error {
					playlist.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	playlistCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids         []string
	title       string
	description string
	tags        []string
	language    string
	privacy     string
	output      string
	yes         bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringSliceVarP(&o.tags, "tags", "a", nil, tagsUsage)
	c.Flags().StringVarP(&o.language, "language", "l", "", languageUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *playlist.Playlist {
	return playlist.NewPlaylist(
		playlist.WithIds(o.ids),
		playlist.WithTitle(o.title),
		utils.OptIfChanged(f, "description", &o.description, playlist.WithDescription),
		utils.OptIfChanged(f, "tags", o.tags, playlist.WithTags),
		utils.OptIfChanged(f, "language", &o.language, playlist.WithLanguage),
		playlist.WithPrivacy(o.privacy),
		playlist.WithMaxResults(1),
		playlist.WithOutput(o.output),
	).(*playlist.Playlist)
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

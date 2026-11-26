// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	updateTool    = "playlistItem-update"
	updateIdUsage = "ID of the playlist item to update"
	updateConfirm = "Update playlist item: %s"
	updateShort   = "Update a playlist item"
	updateLong    = "Update a playlist item. Use this tool to update a playlist item."
	updateExample = `  # Update playlist item title
  yutu playlistItem update --id abc123 --title 'Updated Title'
  # Update playlist item description and privacy
  yutu playlistItem update --id abc123 --description 'New description' --privacy public
  # Update playlist item privacy with JSON output
  yutu playlistItem update --id abc123 --privacy private --output json`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
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
				func(input playlistItem.PlaylistItem) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input playlistItem.PlaylistItem, w io.Writer) error {
					playlistItem.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	playlistItemCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids         []string
	title       string
	description string
	privacy     string
	oboco       string
	output      string
	yes         bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *playlistItem.PlaylistItem {
	return playlistItem.NewPlaylistItem(
		playlistItem.WithIds(o.ids),
		playlistItem.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, playlistItem.WithDescription,
		),
		playlistItem.WithPrivacy(o.privacy),
		playlistItem.WithOBOCO(o.oboco),
		playlistItem.WithMaxResults(1),
		playlistItem.WithOutput(o.output),
	).(*playlistItem.PlaylistItem)
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

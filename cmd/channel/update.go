// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	updateTool    = "channel-update"
	updateIdUsage = "ID of the channel to update"
	updateConfirm = "Update channel: %s"
	updateShort   = "Update a channel"
	updateLong    = "Update a channel. Use this tool to update a channel."
	updateExample = `  # Update channel description
  yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --description 'New description'
  # Update channel title and country
  yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --title 'New Title' --country US
  # Update channel default language
  yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --defaultLanguage en
  # Update channel keywords
  yutu channel update --id UC_x5XG1OV2P6uZZ5FSM9Ttw --keywords 'AI agents,MCP'`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"country":          {Type: "string", Description: countryUsage},
		"default_language": {Type: "string", Description: dlUsage},
		"description":      {Type: "string", Description: descUsage},
		"keywords": {
			Type: "array", Description: keywordsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"title": {Type: "string", Description: titleUsage},
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
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input channel.Channel) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input channel.Channel, w io.Writer) error {
					channel.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	channelCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids             []string
	country         string
	defaultLanguage string
	description     string
	keywords        []string
	title           string
	output          string
	yes             bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().StringVarP(&o.country, "country", "c", "", countryUsage)
	c.Flags().StringVarP(
		&o.defaultLanguage, "defaultLanguage", "l", "", dlUsage,
	)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringSliceVarP(&o.keywords, "keywords", "k", nil, keywordsUsage)
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *channel.Channel {
	return channel.NewChannel(
		channel.WithIds(o.ids),
		utils.OptIfChanged(f, "country", &o.country, channel.WithCountry),
		utils.OptIfChanged(
			f, "defaultLanguage", &o.defaultLanguage, channel.WithDefaultLanguage,
		),
		utils.OptIfChanged(f, "description", &o.description, channel.WithDescription),
		utils.OptIfChanged(f, "keywords", o.keywords, channel.WithKeywords),
		channel.WithTitle(o.title),
		channel.WithMaxResults(1),
		channel.WithOutput(o.output),
	).(*channel.Channel)
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

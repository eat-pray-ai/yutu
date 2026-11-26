// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"context"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/i18nLanguage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort = "List i18n languages"
	listLong  = `List i18n languages. Use this tool to list available internationalization languages.`
)

func init() {
	cmd.Server.AddResource(
		&mcp.Resource{
			URI:         hlURI,
			Name:        hlName,
			Description: hlDesc,
			MIMEType:    pkg.JsonMIME,
			Annotations: &mcp.Annotations{
				Audience: []mcp.Role{"user", "assistant"}, Priority: 0.51,
			},
		}, hlHandler,
	)
	cmd.Server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        langName,
			Description: listLong,
			MIMEType:    pkg.JsonMIME,
			URITemplate: langURI,
			Annotations: &mcp.Annotations{
				Audience: []mcp.Role{"user", "assistant"}, Priority: 0.51,
			},
		}, langsHandler,
	)
	i18nLanguageCmd.AddCommand(newListCmd())
}

type listOpts struct {
	hl     string
	parts  []string
	output string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", defaultParts, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *i18nLanguage.I18nLanguage {
	return i18nLanguage.NewI18nLanguage(
		i18nLanguage.WithHl(o.hl),
		i18nLanguage.WithParts(o.parts),
		i18nLanguage.WithOutput(o.output),
	).(*i18nLanguage.I18nLanguage)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:   "list",
		Short: listShort,
		Long:  listLong,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

var hlHandler = cobramcp.GenResourceHandler(
	hlName, pkg.JsonMIME, func(ctx context.Context, req *mcp.ReadResourceRequest, w io.Writer) error {
		input := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithContext(ctx),
			i18nLanguage.WithParts(defaultParts),
			i18nLanguage.WithOutput("json"),
		)
		return input.List(w)
	},
)

var langsHandler = cobramcp.GenResourceHandler(
	langName, pkg.JsonMIME, func(ctx context.Context, req *mcp.ReadResourceRequest, w io.Writer) error {
		hl := utils.ExtractHl(req.Params.URI)
		input := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithContext(ctx),
			i18nLanguage.WithHl(hl),
			i18nLanguage.WithParts(defaultParts),
			i18nLanguage.WithOutput("json"),
		)
		return input.List(w)
	},
)

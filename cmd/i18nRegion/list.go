// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"context"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/i18nRegion"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort = "List i18n regions"
	listLong  = `List i18n regions. Use this tool to list available internationalization regions.`
)

func init() {
	cmd.Server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        regionName,
			Description: listLong,
			MIMEType:    pkg.JsonMIME,
			URITemplate: regionURI,
			Annotations: &mcp.Annotations{
				Audience: []mcp.Role{"user", "assistant"}, Priority: 0.51,
			},
		}, regionsHandler,
	)
	i18nRegionCmd.AddCommand(newListCmd())
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

func (o *listOpts) build() *i18nRegion.I18nRegion {
	return i18nRegion.NewI18nRegion(
		i18nRegion.WithHl(o.hl),
		i18nRegion.WithParts(o.parts),
		i18nRegion.WithOutput(o.output),
	).(*i18nRegion.I18nRegion)
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

var regionsHandler = cobramcp.GenResourceHandler(
	regionName, pkg.JsonMIME,
	func(ctx context.Context, req *mcp.ReadResourceRequest, w io.Writer) error {
		hl := utils.ExtractHl(req.Params.URI)
		input := i18nRegion.NewI18nRegion(
			i18nRegion.WithContext(ctx),
			i18nRegion.WithHl(hl),
			i18nRegion.WithParts(defaultParts),
			i18nRegion.WithOutput("json"),
		)
		return input.List(w)
	},
)

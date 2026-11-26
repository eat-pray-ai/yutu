// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"context"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/videoCategory"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        vcName,
			Title:       short,
			Description: long,
			MIMEType:    pkg.JsonMIME,
			URITemplate: vcURI,
			Annotations: &mcp.Annotations{
				Audience: []mcp.Role{"user", "assistant"}, Priority: 0.51,
			},
		}, categoriesHandler,
	)
	videoCategoryCmd.AddCommand(newListCmd())
}

const (
	listShort = "List video categories"
	listLong  = `List video categories. Use this tool to list available video categories.`
)

type listOpts struct {
	ids        []string
	hl         string
	regionCode string
	parts      []string
	output     string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, idsUsage)
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().StringVarP(&o.regionCode, "regionCode", "r", "US", rcUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *videoCategory.VideoCategory {
	return videoCategory.NewVideoCategory(
		videoCategory.WithIds(o.ids),
		videoCategory.WithHl(o.hl),
		videoCategory.WithRegionCode(o.regionCode),
		videoCategory.WithParts(o.parts),
		videoCategory.WithOutput(o.output),
	).(*videoCategory.VideoCategory)
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

var categoriesHandler = cobramcp.GenResourceHandler(
	vcName, pkg.JsonMIME, func(ctx context.Context, req *mcp.ReadResourceRequest, w io.Writer) error {
		hl := utils.ExtractHl(req.Params.URI)
		vc := videoCategory.NewVideoCategory(
			videoCategory.WithContext(ctx),
			videoCategory.WithHl(hl),
			videoCategory.WithRegionCode("US"),
			videoCategory.WithParts(defaultParts),
			videoCategory.WithOutput("json"),
		)
		return vc.List(w)
	},
)

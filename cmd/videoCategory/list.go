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
	videoCategoryCmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", nil, idsUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "US", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringP("output", "o", "table", pkg.TableUsage)
}

const (
	listShort = "List video categories"
	listLong  = `List video categories. Use this tool to list available video categories.`
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	RunE: func(c *cobra.Command, _ []string) error {
		output, _ := c.Flags().GetString("output")
		input := videoCategory.NewVideoCategory(
			videoCategory.WithIds(ids),
			videoCategory.WithHl(hl),
			videoCategory.WithRegionCode(regionCode),
			videoCategory.WithParts(parts),
			videoCategory.WithOutput(output),
		)
		return input.List(c.OutOrStdout())
	},
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

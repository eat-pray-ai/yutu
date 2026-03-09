// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
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
	i18nRegionCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", defaultParts, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := i18nRegion.NewI18nRegion(
			i18nRegion.WithHl(hl),
			i18nRegion.WithParts(parts),
			i18nRegion.WithOutput(output),
			i18nRegion.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

var regionsHandler = cobramcp.GenResourceHandler(
	regionName, pkg.JsonMIME,
	func(req *mcp.ReadResourceRequest, w io.Writer) error {
		hl := utils.ExtractHl(req.Params.URI)
		input := i18nRegion.NewI18nRegion(
			i18nRegion.WithHl(hl),
			i18nRegion.WithParts(defaultParts),
			i18nRegion.WithOutput("json"),
		)
		return input.List(w)
	},
)

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/i18nLanguage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
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
			Description: long,
			MIMEType:    pkg.JsonMIME,
			URITemplate: langURI,
			Annotations: &mcp.Annotations{
				Audience: []mcp.Role{"user", "assistant"}, Priority: 0.51,
			},
		}, langsHandler,
	)
	i18nLanguageCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", defaultParts, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		input := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithHl(hl),
			i18nLanguage.WithParts(parts),
			i18nLanguage.WithOutput(output),
			i18nLanguage.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

var hlHandler = cmd.GenResourceHandler(
	hlName, func(req *mcp.ReadResourceRequest, w io.Writer) error {
		input := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithParts(defaultParts),
			i18nLanguage.WithOutput("json"),
			i18nLanguage.WithJsonpath("$.*.snippet.hl"),
		)
		return input.List(w)
	},
)

var langsHandler = cmd.GenResourceHandler(
	langName, func(req *mcp.ReadResourceRequest, w io.Writer) error {
		hl := utils.ExtractHl(req.Params.URI)
		input := i18nLanguage.NewI18nLanguage(
			i18nLanguage.WithHl(hl),
			i18nLanguage.WithParts(defaultParts),
			i18nLanguage.WithOutput("json"),
		)
		return input.List(w)
	},
)

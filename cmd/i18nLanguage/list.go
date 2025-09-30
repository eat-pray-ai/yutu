// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"bytes"
	"context"
	"io"
	"log/slog"

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
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func hlHandler(
	ctx context.Context, request *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	parts = defaultParts
	output = "json"
	jpath = "$.*.snippet.hl"

	slog.InfoContext(ctx, "i18nLanguage hl list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nLanguage hl list failed", "error", err, "uri",
			request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nLanguage hl list completed successfully",
		"resultSize", writer.Len(),
	)

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI: hlURI, MIMEType: pkg.JsonMIME, Text: writer.String(),
			},
		},
	}, nil
}

func langsHandler(
	ctx context.Context, request *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	parts = defaultParts
	hl = utils.ExtractHl(request.Params.URI)
	output = "json"

	slog.InfoContext(ctx, "i18nLanguage list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nLanguage list failed", "error", err, "uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nLanguage list completed successfully",
		"resultSize", writer.Len(),
	)

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI: request.Params.URI, MIMEType: pkg.JsonMIME, Text: writer.String(),
			},
		},
	}, nil
}

func list(writer io.Writer) error {
	i := i18nLanguage.NewI18nLanguage(
		i18nLanguage.WithHl(hl), i18nLanguage.WithService(nil),
	)

	return i.List(parts, output, jpath, writer)
}

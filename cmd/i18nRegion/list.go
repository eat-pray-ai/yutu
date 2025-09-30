// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/i18nRegion"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			Name:        regionName,
			Description: long,
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

func regionsHandler(
	ctx context.Context, request *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	parts = defaultParts
	hl = utils.ExtractHl(request.Params.URI)
	output = "json"

	slog.InfoContext(ctx, "i18nRegion list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "i18nRegion list failed", "error", err, "uri", request.Params.URI,
		)
		return nil, err
	}

	slog.InfoContext(
		ctx, "i18nRegion list completed successfully",
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
	i := i18nRegion.NewI18nRegion(
		i18nRegion.WithHl(hl), i18nRegion.WithService(nil),
	)

	return i.List(parts, output, jpath, writer)
}

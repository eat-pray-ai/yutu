// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

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
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
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
	ctx context.Context, req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: regionName, MinInterval: time.Second,
			},
		),
	)

	parts = defaultParts
	hl = utils.ExtractHl(req.Params.URI)
	output = "json"

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "uri", req.Params.URI)
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{URI: req.Params.URI, MIMEType: pkg.JsonMIME, Text: writer.String()},
		},
	}, nil
}

func list(writer io.Writer) error {
	i := i18nRegion.NewI18nRegion(
		i18nRegion.WithHl(hl), i18nRegion.WithService(nil),
	)

	return i.List(parts, output, jsonpath, writer)
}

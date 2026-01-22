// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

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
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "US", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
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

func categoriesHandler(
	ctx context.Context, req *mcp.ReadResourceRequest,
) (*mcp.ReadResourceResult, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: vcName, MinInterval: time.Second},
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
	vc := videoCategory.NewVideoCategory(
		videoCategory.WithIDs(ids),
		videoCategory.WithHl(hl),
		videoCategory.WithRegionCode(regionCode),
		videoCategory.WithService(nil),
	)

	return vc.List(parts, output, jsonpath, writer)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GenResourceHandler(
	name string, op func(*mcp.ReadResourceRequest, io.Writer) error,
) mcp.ResourceHandler {
	return func(
		ctx context.Context, req *mcp.ReadResourceRequest,
	) (*mcp.ReadResourceResult, error) {
		logger := slog.New(
			mcp.NewLoggingHandler(
				req.Session,
				&mcp.LoggingHandlerOptions{
					LoggerName: name, MinInterval: time.Second,
				},
			),
		)

		var writer bytes.Buffer
		err := op(req, &writer)
		if err != nil {
			logger.ErrorContext(ctx, err.Error(), "uri", req.Params.URI)
			slog.ErrorContext(ctx, err.Error(), "uri", req.Params.URI)
			return nil, err
		}

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{URI: req.Params.URI, MIMEType: pkg.JsonMIME, Text: writer.String()},
			},
		}, nil
	}
}

func GenToolHandler[T any](
	toolName string, op func(T, io.Writer) error,
) mcp.ToolHandlerFor[T, any] {
	return func(
		ctx context.Context, req *mcp.CallToolRequest, input T,
	) (*mcp.CallToolResult, any, error) {
		logger := slog.New(
			mcp.NewLoggingHandler(
				req.Session,
				&mcp.LoggingHandlerOptions{
					LoggerName: toolName, MinInterval: time.Second,
				},
			),
		)

		var writer bytes.Buffer
		err := op(input, &writer)
		if err != nil {
			logger.ErrorContext(ctx, err.Error(), "input", input)
			slog.ErrorContext(ctx, err.Error(), "tool", toolName, "input", input)
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
	}
}

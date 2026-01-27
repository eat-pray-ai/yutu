// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
			return nil, nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
	}
}

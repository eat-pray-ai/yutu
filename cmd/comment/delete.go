// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool  = "comment-delete"
	deleteShort = "Delete YouTube comments"
	deleteLong  = "Delete YouTube comments by ids"
)

type deleteIn struct {
	IDs []string `json:"ids"`
}

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: idsUsage,
			Default:     json.RawMessage(`[]`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	commentCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func deleteHandler(
	ctx context.Context, req *mcp.CallToolRequest, input deleteIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: deleteTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.IDs
	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func del(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithService(nil),
	)

	return c.Delete(writer)
}

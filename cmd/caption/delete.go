// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "caption-delete"
	deleteShort    = "Delete captions"
	deleteLong     = "Delete captions of a video by ids"
	deleteIdsUsage = "IDs of the captions to delete"
)

type deleteIn struct {
	Ids                    []string `json:"ids"`
	OnBehalfOf             string   `json:"onBehalfOf"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
}

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"onBehalfOf":             {Type: "string"},
		"onBehalfOfContentOwner": {Type: "string"},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	captionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := &deleteIn{
			Ids:                    ids,
			OnBehalfOf:             onBehalfOf,
			OnBehalfOfContentOwner: onBehalfOfContentOwner,
		}
		if err := input.call(cmd.OutOrStdout()); err != nil {
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

	var writer bytes.Buffer
	if err := input.call(&writer); err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func (d *deleteIn) call(writer io.Writer, opts ...caption.Option) error {
	defaultOpts := []caption.Option{
		caption.WithIDs(d.Ids),
		caption.WithOnBehalfOf(d.OnBehalfOf),
		caption.WithOnBehalfOfContentOwner(d.OnBehalfOfContentOwner),
		caption.WithService(nil),
	}
	defaultOpts = append(defaultOpts, opts...)

	c := caption.NewCation(defaultOpts...)

	return c.Delete(writer)
}

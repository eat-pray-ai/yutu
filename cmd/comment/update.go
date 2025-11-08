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
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "comment-update"
	updateShort   = "Update a comment"
	updateLong    = "Update a comment on a video"
	updateIdUsage = "ID of the comment"
)

type updateIn struct {
	IDs          []string `json:"ids"`
	CanRate      string   `json:"canRate"`
	TextOriginal string   `json:"textOriginal"`
	ViewerRating string   `json:"viewerRating"`
	Output       string   `json:"output"`
	Jsonpath     string   `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "canRate", "textOriginal", "viewerRating",
		"output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
			Default:     json.RawMessage(`[]`),
		},
		"canRate": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: crUsage, Default: json.RawMessage(`""`),
		},
		"textOriginal": {
			Type: "string", Description: toUsage,
			Default: json.RawMessage(`""`),
		},
		"viewerRating": {
			Type: "string", Enum: []any{"none", "like", "dislike", ""},
			Description: vrUsage, Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	commentCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	updateCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	updateCmd.Flags().StringVarP(
		&viewerRating, "viewerRating", "r", "", vrUsage,
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, req *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: updateTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.IDs
	canRate = utils.BoolPtr(input.CanRate)
	textOriginal = input.TextOriginal
	viewerRating = input.ViewerRating
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func update(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithCanRate(canRate),
		comment.WithTextOriginal(textOriginal),
		comment.WithViewerRating(viewerRating),
		comment.WithMaxResults(1),
		comment.WithService(nil),
	)

	return c.Update(output, jpath, writer)
}

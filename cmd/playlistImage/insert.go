// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool  = "playlistImage-insert"
	insertShort = "Insert a YouTube playlist image"
	insertLong  = "Insert a YouTube playlist image for a given playlist id"
)

type insertIn struct {
	File       string `json:"file"`
	PlaylistId string `json:"playlistId"`
	Type       string `json:"type"`
	Height     int64  `json:"height"`
	Width      int64  `json:"width"`
	Output     string `json:"output"`
	Jsonpath   string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "playlistId"},
	Properties: map[string]*jsonschema.Schema{
		"file": {
			Type: "string", Description: fileUsage,
			Default: json.RawMessage(`""`),
		},
		"playlistId": {
			Type: "string", Description: pidUsage,
			Default: json.RawMessage(`""`),
		},
		"type": {
			Type: "string", Description: typeUsage,
			Default: json.RawMessage(`""`),
		},
		"height": {
			Type: "number", Description: heightUsage,
			Default: json.RawMessage("0"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"width": {
			Type: "number", Description: widthUsage,
			Default: json.RawMessage("0"),
			Minimum: jsonschema.Ptr(float64(0)),
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	file = input.File
	playlistId = input.PlaylistId
	type_ = input.Type
	height = input.Height
	width = input.Width
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithFile(file),
		playlistImage.WithPlaylistID(playlistId),
		playlistImage.WithType(type_),
		playlistImage.WithHeight(height),
		playlistImage.WithWidth(width),
		playlistImage.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

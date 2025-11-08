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
	updateTool  = "playlistImage-update"
	updateShort = "Update a playlist image"
	updateLong  = "Update a playlist image for a given playlist id"
)

type updateIn struct {
	PlaylistId string `json:"playlistId"`
	Type       string `json:"type"`
	Height     int64  `json:"height"`
	Width      int64  `json:"width"`
	Output     string `json:"output"`
	Jsonpath   string `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"playlistId", "type", "height", "width", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	playlistImageCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	updateCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	updateCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	updateCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("playlistId")
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

	playlistId = input.PlaylistId
	type_ = input.Type
	height = input.Height
	width = input.Width
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
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithPlaylistID(playlistId),
		playlistImage.WithType(type_),
		playlistImage.WithHeight(height),
		playlistImage.WithWidth(width),
		playlistImage.WithMaxResults(1),
		playlistImage.WithService(nil),
	)

	return pi.Update(output, jpath, writer)
}

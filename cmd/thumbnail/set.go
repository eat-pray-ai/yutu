// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	setTool = "thumbnail-set"
)

type setIn struct {
	File     string `json:"file"`
	VideoId  string `json:"videoId"`
	Output   string `json:"output"`
	Jsonpath string `json:"jsonpath"`
}

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "videoId"},
	Properties: map[string]*jsonschema.Schema{
		"file": {
			Type: "string", Description: fileUsage,
			Default: json.RawMessage(`""`),
		},
		"videoId": {
			Type: "string", Description: vidUsage,
			Default: json.RawMessage(`""`),
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
			Name: setTool, Title: short, Description: long,
			InputSchema: setInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, setHandler,
	)
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	setCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setHandler(
	ctx context.Context, req *mcp.CallToolRequest, input setIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: setTool, MinInterval: time.Second},
		),
	)

	file = input.File
	videoId = input.VideoId
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := set(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func set(writer io.Writer) error {
	t := thumbnail.NewThumbnail(
		thumbnail.WithFile(file),
		thumbnail.WithVideoId(videoId),
		thumbnail.WithService(nil),
	)

	return t.Set(output, jsonpath, writer)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"bytes"
	"context"
	"encoding/json"
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

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":     {Type: "string", Description: fileUsage},
		"video_id": {Type: "string", Description: vidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
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
		input := thumbnail.NewThumbnail(
			thumbnail.WithFile(file),
			thumbnail.WithVideoId(videoId),
			thumbnail.WithOutput(output),
			thumbnail.WithJsonpath(jsonpath),
		)
		err := input.Set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setHandler(
	ctx context.Context, req *mcp.CallToolRequest, input thumbnail.Thumbnail,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: setTool, MinInterval: time.Second},
		),
	)

	var writer bytes.Buffer
	err := input.Set(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

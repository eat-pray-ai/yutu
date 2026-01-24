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
	downloadTool    = "caption-download"
	downloadShort   = "Download caption"
	downloadLong    = "Download caption from a video"
	downloadIdUsage = "ID of the caption to download"
)

type downloadIn struct {
	Ids                    []string `json:"ids"`
	File                   string   `json:"file"`
	Tfmt                   string   `json:"tfmt"`
	Tlang                  string   `json:"tlang"`
	OnBehalfOf             string   `json:"onBehalfOf"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
}

var downloadInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "file"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: downloadIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"file": {Type: "string", Description: fileUsage},
		"tfmt": {
			Type: "string", Description: tfmtUsage,
			Enum: []any{"sbv", "srt", "vtt", ""},
		},
		"tlang":                  {Type: "string", Description: tlangUsage},
		"onBehalfOf":             {Type: "string"},
		"onBehalfOfContentOwner": {Type: "string"},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: downloadTool, Title: downloadShort, Description: downloadLong,
			InputSchema: downloadInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, downloadHandler,
	)
	captionCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringSliceVarP(
		&ids, "id", "i", []string{}, downloadIdUsage,
	)
	downloadCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", tfmtUsage)
	downloadCmd.Flags().StringVarP(&tlang, "tlang", "l", "", tlangUsage)
	downloadCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)

	_ = downloadCmd.MarkFlagRequired("id")
	_ = downloadCmd.MarkFlagRequired("file")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: downloadShort,
	Long:  downloadLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := &downloadIn{
			Ids:                    ids,
			File:                   file,
			Tfmt:                   tfmt,
			Tlang:                  tlang,
			OnBehalfOf:             onBehalfOf,
			OnBehalfOfContentOwner: onBehalfOfContentOwner,
		}
		if err := input.call(cmd.OutOrStdout()); err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func downloadHandler(
	ctx context.Context, req *mcp.CallToolRequest, input downloadIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: downloadTool, MinInterval: time.Second,
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

func (d *downloadIn) call(writer io.Writer, opts ...caption.Option) error {
	defaultOpts := []caption.Option{
		caption.WithIDs(d.Ids),
		caption.WithFile(d.File),
		caption.WithTfmt(d.Tfmt),
		caption.WithTlang(d.Tlang),
		caption.WithOnBehalfOf(d.OnBehalfOf),
		caption.WithOnBehalfOfContentOwner(d.OnBehalfOfContentOwner),
		caption.WithService(nil),
	}
	defaultOpts = append(defaultOpts, opts...)

	c := caption.NewCation(defaultOpts...)

	return c.Download(writer)
}

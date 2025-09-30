// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
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
	Type: "object",
	Required: []string{
		"ids", "file", "tfmt", "tlang", "onBehalfOf", "onBehalfOfContentOwner",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: downloadIdUsage,
			Default:     json.RawMessage(`[]`),
		},
		"file": {
			Type:        "string",
			Description: fileUsage,
			Default:     json.RawMessage(`""`),
		},
		"tfmt": {
			Type:        "string",
			Enum:        []any{"sbv", "srt", "vtt", ""},
			Description: tfmtUsage,
			Default:     json.RawMessage(`""`),
		},
		"tlang": {
			Type:        "string",
			Description: tlangUsage,
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOf": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "caption-download", Title: downloadShort, Description: downloadLong,
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
		err := download(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func downloadHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input downloadIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	file = input.File
	tfmt = input.Tfmt
	tlang = input.Tlang
	onBehalfOf = input.OnBehalfOf
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	slog.InfoContext(ctx, "caption download started")

	var writer bytes.Buffer
	err := download(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "caption download failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "caption download completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func download(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithIDs(ids),
		caption.WithFile(file),
		caption.WithTfmt(tfmt),
		caption.WithTlang(tlang),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithService(nil),
	)

	return c.Download(writer)
}

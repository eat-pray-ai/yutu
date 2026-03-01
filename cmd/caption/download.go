// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	downloadTool    = "caption-download"
	downloadShort   = "Download caption"
	downloadLong    = "Download caption from a video\n\nExamples:\n  yutu caption download --id abc123 --file subtitle.srt\n  yutu caption download --id abc123 --file subtitle.vtt --tfmt vtt\n  yutu caption download --id abc123 --file subtitle.srt --tlang fr"
	downloadIdUsage = "ID of the caption to download"
)

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
		"tlang":                      {Type: "string", Description: tlangUsage},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: downloadTool, Title: downloadShort, Description: downloadLong,
			InputSchema: downloadInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			downloadTool, func(input caption.Caption, writer io.Writer) error {
				return input.Download(writer)
			},
		),
	)
	captionCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringSliceVarP(
		&ids, "id", "i", []string{}, downloadIdUsage,
	)
	downloadCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	downloadCmd.Flags().StringVarP(&tfmt, "tfmt", "t", "", tfmtUsage)
	downloadCmd.Flags().StringVarP(&tlang, "tlang", "l", "", tlangUsage)
	downloadCmd.Flags().StringVarP(
		&onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage,
	)
	downloadCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)

	_ = downloadCmd.MarkFlagRequired("id")
	_ = downloadCmd.MarkFlagRequired("file")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: downloadShort,
	Long:  downloadLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := caption.NewCaption(
			caption.WithIds(ids),
			caption.WithFile(file),
			caption.WithTfmt(tfmt),
			caption.WithTlang(tlang),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		utils.HandleCmdError(input.Download(cmd.OutOrStdout()), cmd)
	},
}

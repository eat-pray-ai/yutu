// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"encoding/json"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	setTool    = "thumbnail-set"
	setShort   = "Set a thumbnail for a video"
	setLong    = "Set a thumbnail for a video. Use this tool when you need to set a thumbnail for a video."
	setExample = `yutu thumbnail set --file image.jpg --videoId dQw4w9WgXcQ
yutu thumbnail set --file image.png --videoId dQw4w9WgXcQ --output json`
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
			Name: setTool, Title: setShort, Description: setLong,
			InputSchema: setInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			setTool, func(input thumbnail.Thumbnail, writer io.Writer) error {
				return input.Set(writer)
			},
		),
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
	Use:     "set",
	Short:   setShort,
	Long:    setLong,
	Example: setExample,
	Run: func(cmd *cobra.Command, args []string) {
		input := thumbnail.NewThumbnail(
			thumbnail.WithFile(file),
			thumbnail.WithVideoId(videoId),
			thumbnail.WithOutput(output),
			thumbnail.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Set(cmd.OutOrStdout()), cmd)
	},
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"encoding/json/jsontext"
	"fmt"
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
	setConfirm = "Set thumbnail for video: %s"
	setShort   = "Set a thumbnail for a video"
	setLong    = "Set a thumbnail for a video. Use this tool to set a thumbnail for a video."
	setExample = `# Set a thumbnail for a video
yutu thumbnail set --file image.jpg --videoId dQw4w9WgXcQ
# Set a thumbnail with JSON output
yutu thumbnail set --file image.png --videoId dQw4w9WgXcQ --output json`
)

var setInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "video_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":     {Type: "string", Description: fileUsage},
		"video_id": {Type: "string", Description: vidUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"yaml"`),
		},
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
		}, cobramcp.GenToolHandlerWithMRTR(
			setTool, cobramcp.ConfirmThen(
				func(input thumbnail.Thumbnail) string {
					return fmt.Sprintf(setConfirm, input.VideoId)
				},
				func(input thumbnail.Thumbnail, w io.Writer) error {
					return input.Set(w)
				},
			),
		),
	)
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	setCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

var setCmd = &cobra.Command{
	Use:     "set",
	Short:   setShort,
	Long:    setLong,
	Example: setExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(setConfirm, videoId))
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := thumbnail.NewThumbnail(
			thumbnail.WithFile(file),
			thumbnail.WithVideoId(videoId),
			thumbnail.WithOutput(output),
		)
		utils.HandleCmdError(input.Set(c.OutOrStdout()), c)
	},
}

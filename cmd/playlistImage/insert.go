// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "playlistImage-insert"
	insertConfirm = "Upload playlist image: %s"
	insertShort   = "Upload a playlist image"
	insertLong    = "Upload a playlist image. Use this tool to upload a playlist image."
	insertExample = `# Insert a playlist cover image
yutu playlistImage insert --file cover.jpg --playlistId PLxxx
# Insert a hero image
yutu playlistImage insert --file cover.png --playlistId PLxxx --type hero
# Insert an image with custom dimensions
yutu playlistImage insert --file cover.jpg --playlistId PLxxx --width 2048 --height 1152`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "playlist_id"},
	Properties: map[string]*jsonschema.Schema{
		"file":        {Type: "string", Description: fileUsage},
		"playlist_id": {Type: "string", Description: pidUsage},
		"type":        {Type: "string", Description: typeUsage},
		"height": {
			Type: "number", Description: heightUsage,
			Minimum: new(float64(0)),
		},
		"width": {
			Type: "number", Description: widthUsage,
			Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"yaml"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input playlistImage.PlaylistImage) string {
					return fmt.Sprintf(insertConfirm, input.File)
				},
				func(input playlistImage.PlaylistImage, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		pkg.OBOCOCUsage,
	)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(insertConfirm, file))
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithFile(file),
			playlistImage.WithPlaylistId(playlistId),
			playlistImage.WithType(type_),
			playlistImage.WithHeight(height),
			playlistImage.WithWidth(width),
			playlistImage.WithOutput(output),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		)
		err := pi.Insert(c.OutOrStdout())
		utils.HandleCmdError(err, c)
	},
}

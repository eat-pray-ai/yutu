// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"encoding/json"
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
	insertTool  = "playlistImage-insert"
	insertShort = "Insert a YouTube playlist image"
	insertLong  = "Insert a YouTube playlist image for a given playlist id\n\nExamples:\n  yutu playlistImage insert --file cover.jpg --playlistId PLxxxxxxxxxxxxxxxxxx\n  yutu playlistImage insert --file cover.png --playlistId PLxxx --type hero\n  yutu playlistImage insert --file cover.jpg --playlistId PLxxx --width 2048 --height 1152"
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
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
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
		}, cobramcp.GenToolHandler(
			insertTool,
			func(input playlistImage.PlaylistImage, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		pkg.OBOCOCUsage,
	)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithFile(file),
			playlistImage.WithPlaylistId(playlistId),
			playlistImage.WithType(type_),
			playlistImage.WithHeight(height),
			playlistImage.WithWidth(width),
			playlistImage.WithOutput(output),
			playlistImage.WithJsonpath(jsonpath),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		)
		err := pi.Insert(cmd.OutOrStdout())
		utils.HandleCmdError(err, cmd)
	},
}

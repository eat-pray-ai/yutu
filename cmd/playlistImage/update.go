// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool  = "playlistImage-update"
	updateShort = "Update a playlist image"
	updateLong  = "Update a playlist image for a given playlist id\n\nExamples:\n  yutu playlistImage update --playlistId PLxxxxxxxxxxxxxxxxxx\n  yutu playlistImage update --playlistId PLxxx --type hero --width 2048 --height 1152"
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"playlist_id"},
	Properties: map[string]*jsonschema.Schema{
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
		"on_behalf_of_content_owner":         {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {Type: "string", Description: pkg.OBOCOCUsage},
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			updateTool,
			func(input playlistImage.PlaylistImage, writer io.Writer) error {
				return input.Update(writer)
			},
		),
	)
	playlistImageCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	updateCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	updateCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	updateCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", pkg.OBOCOCUsage,
	)

	_ = updateCmd.MarkFlagRequired("playlistId")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := playlistImage.NewPlaylistImage(
			playlistImage.WithPlaylistId(playlistId),
			playlistImage.WithType(type_),
			playlistImage.WithHeight(height),
			playlistImage.WithWidth(width),
			playlistImage.WithMaxResults(1),
			playlistImage.WithOutput(output),
			playlistImage.WithJsonpath(jsonpath),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		)
		utils.HandleCmdError(input.Update(cmd.OutOrStdout()), cmd)
	},
}

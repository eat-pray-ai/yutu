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
	listTool  = "playlistImage-list"
	listShort = "List YouTube playlist images"
	listLong  = "List YouTube playlist images' info"
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"parent": {Type: "string", Description: parentUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"on_behalf_of_content_owner":         {Type: "string"},
		"on_behalf_of_content_owner_channel": {Type: "string"},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","kind","snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, cmd.GenToolHandler(
			listTool, func(input playlistImage.PlaylistImage, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	playlistImageCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&parent, "parent", "P", "", parentUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "kind", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := playlistImage.NewPlaylistImage(
			playlistImage.WithParent(parent),
			playlistImage.WithMaxResults(maxResults),
			playlistImage.WithParts(parts),
			playlistImage.WithOutput(output),
			playlistImage.WithJsonpath(jsonpath),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			playlistImage.WithService(nil),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "playlist-delete"
	deleteIdsUsage = "IDs of the playlists to delete"
	deleteShort    = "Delete playlists"
	deleteLong     = "Delete playlists. Use this tool to delete playlists by IDs."
	deleteExample  = `# Delete a playlist by ID
yutu playlist delete --ids PLxxxx
# Delete multiple playlists
yutu playlist delete --ids PLxxx1,PLxxx2`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"confirmed":                  {Type: "boolean", Description: pkg.ConfirmedUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(true),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			deleteTool, func(input playlist.Playlist, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Delete(writer)
			},
		),
	)
	playlistCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	deleteCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Long:    deleteLong,
	Example: deleteExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf("Would delete playlist(s): %s", strings.Join(ids, ", "))
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		input := playlist.NewPlaylist(
			playlist.WithIds(ids),
			playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

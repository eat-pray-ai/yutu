// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "playlistItem-update"
	updateIdUsage = "ID of the playlist item to update"
	updateShort   = "Update a playlist item"
	updateLong    = "Update a playlist item. Use this tool to update a playlist item."
	updateExample = `# Update playlist item title
yutu playlistItem update --id abc123 --title 'Updated Title'
# Update playlist item description and privacy
yutu playlistItem update --id abc123 --description 'New description' --privacy public
# Update playlist item privacy with JSON output
yutu playlistItem update --id abc123 --privacy private --output json`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			updateTool, func(input playlistItem.PlaylistItem, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Update(writer)
			},
		),
	)
	playlistItemCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	updateCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf("Would update playlist item: %s", strings.Join(ids, ", "))
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := playlistItem.NewPlaylistItem(
			playlistItem.WithIds(ids),
			playlistItem.WithTitle(title),
			playlistItem.WithDescription(description),
			playlistItem.WithPrivacy(privacy),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}

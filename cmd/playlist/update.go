// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"encoding/json"
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
	updateTool    = "playlist-update"
	updateIdUsage = "ID of the playlist to update"
	updateShort   = "Update a playlist"
	updateLong    = "Update a playlist. Use this tool to update a playlist."
	updateExample = `# Update playlist title
yutu playlist update --id PLxxx --title 'Updated Title'
# Update playlist description and privacy
yutu playlist update --id PLxxx --description 'New description' --privacy public
# Update playlist tags
yutu playlist update --id PLxxx --tags 'music,pop,2024'`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
		},
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"tags": {
			Type: "array", Description: tagsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"language": {Type: "string", Description: languageUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
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
			updateTool, func(input playlist.Playlist, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Update(writer)
			},
		),
	)
	playlistCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
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
		msg := fmt.Sprintf("Would update playlist: %s", strings.Join(ids, ", "))
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		p := playlist.NewPlaylist(
			playlist.WithIds(ids),
			playlist.WithTitle(title),
			playlist.WithDescription(description),
			playlist.WithTags(tags),
			playlist.WithLanguage(language),
			playlist.WithPrivacy(privacy),
			playlist.WithOutput(output),
		)
		utils.HandleCmdError(p.Update(c.OutOrStdout()), c)
	},
}

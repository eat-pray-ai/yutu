// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

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
	insertTool     = "playlistItem-insert"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
	insertConfirm  = "Add playlist item to playlist: %s"
	insertShort    = "Add a playlist item"
	insertLong     = "Add a playlist item. Use this tool to add a playlist item to a playlist."
	insertExample  = `# Add a video to a playlist
yutu playlistItem insert --kind video --playlistId PLxxx --channelId UC_x5X --kVideoId dQw4w9WgXcQ
# Add a video to a playlist with privacy setting
yutu playlistItem insert --kind video --playlistId PLxxx --channelId UC_x5X --kVideoId dQw4w9 --privacy public`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"kind", "playlist_id", "channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"kind": {
			Type: "string", Description: kindUsage,
			Enum: []any{"video", "channel", "playlist"},
		},
		"k_video_id":    {Type: "string", Description: kvidUsage},
		"k_channel_id":  {Type: "string", Description: kcidUsage},
		"k_playlist_id": {Type: "string", Description: kpidUsage},
		"playlist_id":   {Type: "string", Description: insertPidUsage},
		"channel_id":    {Type: "string", Description: cidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
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
				func(input playlistItem.PlaylistItem) string {
					return fmt.Sprintf(insertConfirm, input.PlaylistId)
				},
				func(input playlistItem.PlaylistItem, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&kind, "kind", "k", "", kindUsage)
	insertCmd.Flags().StringVarP(&kVideoId, "kVideoId", "V", "", kvidUsage)
	insertCmd.Flags().StringVarP(&kChannelId, "kChannelId", "C", "", kcidUsage)
	insertCmd.Flags().StringVarP(&kPlaylistId, "kPlaylistId", "Y", "", kpidUsage)
	insertCmd.Flags().StringVarP(
		&playlistId, "playlistId", "y", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = insertCmd.MarkFlagRequired("kind")
	_ = insertCmd.MarkFlagRequired("playlistId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(c, fmt.Sprintf(insertConfirm, playlistId))
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := playlistItem.NewPlaylistItem(
			playlistItem.WithTitle(title),
			playlistItem.WithDescription(description),
			playlistItem.WithKind(kind),
			playlistItem.WithKVideoId(kVideoId),
			playlistItem.WithKChannelId(kChannelId),
			playlistItem.WithKPlaylistId(kPlaylistId),
			playlistItem.WithPlaylistId(playlistId),
			playlistItem.WithPrivacy(privacy),
			playlistItem.WithChannelId(channelId),
			playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistItem.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

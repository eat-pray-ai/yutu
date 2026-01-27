// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool     = "playlistItem-insert"
	insertShort    = "Insert a playlist item into a playlist"
	insertLong     = "Insert a playlist item into a playlist"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"kind", "playlist_id", "channel_id"},
	Properties: map[string]*jsonschema.Schema{
		"title":       {Type: "string", Description: titleUsage},
		"description": {Type: "string", Description: descUsage},
		"kind": {
			Type: "string", Description: kindUsage,
			Enum: []any{"video", "channel", "playlist", ""},
		},
		"k_video_id":    {Type: "string", Description: kvidUsage},
		"k_channel_id":  {Type: "string", Description: kcidUsage},
		"k_playlist_id": {Type: "string", Description: kpidUsage},
		"playlist_id":   {Type: "string", Description: insertPidUsage},
		"channel_id":    {Type: "string", Description: cidUsage},
		"privacy": {
			Type: "string", Description: privacyUsage,
			Enum: []any{"public", "private", "unlisted", ""},
		},
		"on_behalf_of_content_owner": {Type: "string"},
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
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, cmd.GenToolHandler(
			insertTool, func(input playlistItem.PlaylistItem, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
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
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("kind")
	_ = insertCmd.MarkFlagRequired("playlistId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
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
			playlistItem.WithJsonpath(jsonpath),
		)
		err := input.Insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

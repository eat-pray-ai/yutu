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
	"github.com/spf13/pflag"
)

const (
	insertTool     = "playlistItem-insert"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
	insertConfirm  = "Add playlist item to playlist %s"
	insertShort    = "Add a playlist item"
	insertLong     = "Add a playlist item. Use this tool to add a playlist item to a playlist."
	insertExample  = `  # Add a video to a playlist
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
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
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
	playlistItemCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	title       string
	description string
	kind        string
	kVideoId    string
	kChannelId  string
	kPlaylistId string
	playlistId  string
	channelId   string
	privacy     string
	oboco       string
	output      string
	yes         bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(&o.kind, "kind", "k", "", kindUsage)
	c.Flags().StringVarP(&o.kVideoId, "kVideoId", "V", "", kvidUsage)
	c.Flags().StringVarP(&o.kChannelId, "kChannelId", "C", "", kcidUsage)
	c.Flags().StringVarP(&o.kPlaylistId, "kPlaylistId", "Y", "", kpidUsage)
	c.Flags().StringVarP(
		&o.playlistId, "playlistId", "y", "", insertPidUsage,
	)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().StringVarP(&o.privacy, "privacy", "p", "", privacyUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("kind")
	_ = c.MarkFlagRequired("playlistId")
	_ = c.MarkFlagRequired("channelId")
}

func (o *insertOpts) build(f *pflag.FlagSet) *playlistItem.PlaylistItem {
	return playlistItem.NewPlaylistItem(
		playlistItem.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, playlistItem.WithDescription,
		),
		playlistItem.WithKind(o.kind),
		playlistItem.WithKVideoId(o.kVideoId),
		playlistItem.WithKChannelId(o.kChannelId),
		playlistItem.WithKPlaylistId(o.kPlaylistId),
		playlistItem.WithPlaylistId(o.playlistId),
		playlistItem.WithPrivacy(o.privacy),
		playlistItem.WithChannelId(o.channelId),
		playlistItem.WithOBOCO(o.oboco),
		playlistItem.WithOutput(o.output),
	).(*playlistItem.PlaylistItem)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.playlistId))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

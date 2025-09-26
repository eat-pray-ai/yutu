package playlistItem

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertShort    = "Insert a playlist item into a playlist"
	insertLong     = "Insert a playlist item into a playlist"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
)

type insertIn struct {
	Title                  string `json:"title"`
	Description            string `json:"description"`
	Kind                   string `json:"kind"`
	KVideoId               string `json:"kVideoId"`
	KChannelId             string `json:"kChannelId"`
	KPlaylistId            string `json:"kPlaylistId"`
	PlaylistId             string `json:"playlistId"`
	ChannelId              string `json:"channelId"`
	Privacy                string `json:"privacy"`
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"`
	Output                 string `json:"output"`
	Jsonpath               string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"title", "description", "kind", "kVideoId", "kChannelId",
		"kPlaylistId", "playlistId", "channelId", "privacy",
		"onBehalfOfContentOwner", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"kind": {
			Type: "string", Enum: []any{"video", "channel", "playlist", ""},
			Description: kindUsage, Default: json.RawMessage(`""`),
		},
		"kVideoId": {
			Type: "string", Description: kvidUsage,
			Default: json.RawMessage(`""`),
		},
		"kChannelId": {
			Type: "string", Description: kcidUsage,
			Default: json.RawMessage(`""`),
		},
		"kPlaylistId": {
			Type: "string", Description: kpidUsage,
			Default: json.RawMessage(`""`),
		},
		"playlistId": {
			Type: "string", Description: insertPidUsage,
			Default: json.RawMessage(`""`),
		},
		"channelId": {
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "playlistItem-insert", Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
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
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("kind")
	_ = insertCmd.MarkFlagRequired("playlistId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	title = input.Title
	description = input.Description
	kind = input.Kind
	kVideoId = input.KVideoId
	kChannelId = input.KChannelId
	kPlaylistId = input.KPlaylistId
	playlistId = input.PlaylistId
	channelId = input.ChannelId
	privacy = input.Privacy
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "playlistItem insert started")

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlistItem insert failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "playlistItem insert completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func insert(writer io.Writer) error {
	pi := playlistItem.NewPlaylistItem(
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
		playlistItem.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

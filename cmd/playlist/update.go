package playlist

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateShort   = "Update an existing playlist"
	updateLong    = "Update an existing playlist, with the specified title, description, tags, etc"
	updateIdUsage = "ID of the playlist to update"
)

type updateIn struct {
	Ids         []string `json:"ids"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Language    string   `json:"language"`
	Privacy     string   `json:"privacy"`
	Output      string   `json:"output"`
	Jsonpath    string   `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "title", "description", "tags",
		"language", "privacy", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
			Default:     json.RawMessage(`[]`),
		},
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"tags": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: tagsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"language": {
			Type: "string", Description: languageUsage,
			Default: json.RawMessage(`""`),
		},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage, Default: json.RawMessage(`""`),
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
			Name: "playlist-update", Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	playlistCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	title = input.Title
	description = input.Description
	tags = input.Tags
	language = input.Language
	privacy = input.Privacy
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "playlist update started")

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlist update failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "playlist update completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func update(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithIDs(ids),
		playlist.WithTitle(title),
		playlist.WithDescription(description),
		playlist.WithTags(tags),
		playlist.WithLanguage(language),
		playlist.WithPrivacy(privacy),
		playlist.WithMaxResults(1),
		playlist.WithService(nil),
	)

	return p.Update(output, jpath, writer)
}

package playlist

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete a playlists"
	deleteLong     = "Delete a playlists by ids"
	deleteIdsUsage = "IDs of the playlists to delete"
)

type deleteIn struct {
	Ids                    []string `json:"ids"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
}

var deleteInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "onBehalfOfContentOwner",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: deleteIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "playlist-delete", Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	playlistCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func deleteHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input deleteIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	slog.InfoContext(ctx, "playlist delete started")

	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlist delete failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "playlist delete completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func del(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithIDs(ids),
		playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlist.WithService(nil),
	)

	return p.Delete(writer)
}

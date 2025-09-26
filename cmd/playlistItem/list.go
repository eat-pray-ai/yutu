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
	listShort    = "List playlist items"
	listLong     = "List playlist items' info, such as title, description, etc"
	listIdsUsage = "IDs of the playlist items to list"
	listPidUsage = "Return the playlist items within the given playlist"
)

type listIn struct {
	Ids                    []string `json:"ids"`
	PlaylistId             string   `json:"playlistId"`
	MaxResults             int64    `json:"maxResults"`
	VideoId                string   `json:"videoId"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Parts                  []string `json:"parts"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "playlistId", "maxResults", "videoId",
		"onBehalfOfContentOwner", "parts", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"playlistId": {
			Type: "string", Description: listPidUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"videoId": {
			Type: "string", Description: vidUsage,
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet","status"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
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
			Name: "playlistItem-list", Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
	)
	playlistItemCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&playlistId, "playlistId", "y", "", listPidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func listHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input listIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.Ids
	playlistId = input.PlaylistId
	maxResults = input.MaxResults
	videoId = input.VideoId
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	parts = input.Parts
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "playlistItem list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlistItem list failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "playlistItem list completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithIDs(ids),
		playlistItem.WithPlaylistId(playlistId),
		playlistItem.WithMaxResults(maxResults),
		playlistItem.WithVideoId(videoId),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.List(parts, output, jpath, writer)
}

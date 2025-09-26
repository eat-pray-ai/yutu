package playlistImage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type listIn struct {
	Parent                        string   `json:"parent"`
	MaxResults                    int64    `json:"maxResults"`
	Parts                         []string `json:"parts"`
	Output                        string   `json:"output"`
	Jsonpath                      string   `json:"jsonpath"`
	OnBehalfOfContentOwner        string   `json:"onBehalfOfContentOwner"`
	OnBehalfOfContentOwnerChannel string   `json:"onBehalfOfContentOwnerChannel"`
}

var listInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"parent", "maxResults", "parts", "output", "jsonpath",
		"onBehalfOfContentOwner", "onBehalfOfContentOwnerChannel",
	},
	Properties: map[string]*jsonschema.Schema{
		"parent": {
			Type: "string", Description: parentUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","kind","snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwnerChannel": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
	},
}

const (
	listShort = "List YouTube playlist images"
	listLong  = "List YouTube playlist images' info"
)

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "playlistImage-list", Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
	)
	playlistImageCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&parent, "parent", "P", "", parentUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "kind", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
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
	parent = input.Parent
	maxResults = input.MaxResults
	parts = input.Parts
	output = input.Output
	jpath = input.Jsonpath
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	onBehalfOfContentOwnerChannel = input.OnBehalfOfContentOwnerChannel

	slog.InfoContext(ctx, "playlistImage list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlistImage list failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "playlistImage list completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithParent(parent),
		playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		playlistImage.WithMaxResults(maxResults),
		playlistImage.WithService(nil),
	)

	return pi.List(parts, output, jpath, writer)
}

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List captions"
	listLong     = "List captions of a video"
	listIdsUsage = "IDs of the captions to list"
)

type listIn struct {
	Ids                    []string `json:"ids"`
	VideoId                string   `json:"videoId"`
	OnBehalfOf             string   `json:"onBehalfOf"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Parts                  []string `json:"parts"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "videoId", "onBehalfOf", "onBehalfOfContentOwner",
		"parts", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"videoId": {
			Type:        "string",
			Description: vidUsage,
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOf": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet"]`),
		},
		"output": {
			Type:        "string",
			Enum:        []any{"json", "yaml", "table"},
			Description: pkg.TableUsage,
			Default:     json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type:        "string",
			Description: pkg.JPUsage,
			Default:     json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "caption-list", Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
	)
	captionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
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
	videoId = input.VideoId
	onBehalfOf = input.OnBehalfOf
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	parts = input.Parts
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "caption list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "caption list failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "caption list completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithIDs(ids),
		caption.WithVideoId(videoId),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithService(nil),
	)

	return c.List(parts, output, jpath, writer)
}

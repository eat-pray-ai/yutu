package comment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listShort    = "List YouTube comments"
	listLong     = "List YouTube comments by ids"
	listPidUsage = "Returns replies to the specified comment"
)

type listIn struct {
	IDs        []string `json:"ids"`
	MaxResults int64    `json:"maxResults"`
	ParentId   string   `json:"parentId"`
	TextFormat string   `json:"textFormat"`
	Parts      []string `json:"parts"`
	Output     string   `json:"output"`
	Jsonpath   string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "maxResults", "parentId", "textFormat",
		"parts", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: idsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"parentId": {
			Type: "string", Description: listPidUsage,
			Default: json.RawMessage(`""`),
		},
		"textFormat": {
			Type: "string", Enum: []any{"textFormatUnspecified", "html", "plainText"},
			Description: tfUsage, Default: json.RawMessage(`"html"`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet"]`),
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
			Name: "comment-list", Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
	)
	commentCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, pkg.MRUsage,
	)
	listCmd.Flags().StringVarP(&parentId, "parentId", "P", "", listPidUsage)
	listCmd.Flags().StringVarP(
		&textFormat, "textFormat", "t", "html", tfUsage,
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
	ids = input.IDs
	maxResults = input.MaxResults
	parentId = input.ParentId
	textFormat = input.TextFormat
	parts = input.Parts
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "comment list started")

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "comment list failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "comment list completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithMaxResults(maxResults),
		comment.WithParentId(parentId),
		comment.WithTextFormat(textFormat),
		comment.WithService(nil),
	)

	return c.List(parts, output, jpath, writer)
}

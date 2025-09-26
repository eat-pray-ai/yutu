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
	masShort = "Mark YouTube comments as spam"
	masLong  = "Mark YouTube comments as spam by ids"
)

type markAsSpamIn struct {
	IDs      []string `json:"ids"`
	Output   string   `json:"output"`
	Jsonpath string   `json:"jsonpath"`
}

var markAsSpamInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "output", "jsonpath"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: idsUsage,
			Default:     json.RawMessage(`[]`),
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
			Name: "comment-markAsSpam", Title: masShort, Description: masLong,
			InputSchema: markAsSpamInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, markAsSpamHandler,
	)
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	markAsSpamCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	markAsSpamCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = markAsSpamCmd.MarkFlagRequired("ids")
}

var markAsSpamCmd = &cobra.Command{
	Use:   "markAsSpam",
	Short: masShort,
	Long:  masLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := markAsSpam(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func markAsSpamHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input markAsSpamIn,
) (*mcp.CallToolResult, any, error) {
	ids = input.IDs
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "comment markAsSpam started")

	var writer bytes.Buffer
	err := markAsSpam(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "comment markAsSpam failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "comment markAsSpam completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func markAsSpam(writer io.Writer) error {
	c := comment.NewComment(
		comment.WithIDs(ids),
		comment.WithService(nil),
	)

	return c.MarkAsSpam(output, jpath, writer)
}

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete captions"
	deleteLong     = "Delete captions of a video by ids"
	deleteIdsUsage = "IDs of the captions to delete"
)

type deleteIn struct {
	Ids                    []string `json:"ids"`
	OnBehalfOf             string   `json:"onBehalfOf"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
}

var deleteInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "onBehalfOf", "onBehalfOfContentOwner",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: deleteIdsUsage,
			Default:     json.RawMessage(`[]`),
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
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "caption-delete", Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	captionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
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
	onBehalfOf = input.OnBehalfOf
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner

	slog.InfoContext(ctx, "caption delete started")

	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "caption delete failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "caption delete completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func del(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithIDs(ids),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithService(nil),
	)

	return c.Delete(writer)
}

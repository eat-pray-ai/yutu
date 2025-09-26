package subscription

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteShort    = "Delete a YouTube subscriptions"
	deleteLong     = "Delete a YouTube subscriptions by ids"
	deleteIdsUsage = "IDs of the subscriptions to delete"
)

type deleteIn struct {
	Ids []string `json:"ids"`
}

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: deleteIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "subscription-delete", Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	subscriptionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
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

	slog.InfoContext(ctx, "subscription delete started")

	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "subscription delete failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "subscription delete completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func del(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithIDs(ids), subscription.WithService(nil),
	)

	return s.Delete(writer)
}

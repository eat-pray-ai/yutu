package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
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

const (
	deleteShort    = "Delete a video on YouTube"
	deleteLong     = "Delete a video on YouTube by ids"
	deleteIdsUsage = "IDs of the videos to delete"
)

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "video-delete", Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(true),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, deleteHandler,
	)
	videoCmd.AddCommand(deleteCmd)

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

	slog.InfoContext(ctx, "video delete started")

	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video delete failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "video delete completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func del(writer io.Writer) error {
	v := video.NewVideo(video.WithIDs(ids))

	return v.Delete(writer)
}

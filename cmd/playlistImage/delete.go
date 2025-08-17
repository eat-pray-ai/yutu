package playlistImage

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteShort = "Delete YouTube playlist images"
	deleteLong  = "Delete YouTube playlist images by ids"
)

func init() {
	cmd.MCP.AddTool(deleteTool, deleteHandler)
	playlistImageCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
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

var deleteTool = mcp.NewTool(
	"playlistImage-delete",
	mcp.WithTitleAnnotation(deleteShort),
	mcp.WithDescription(deleteLong),
	mcp.WithDestructiveHintAnnotation(true),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(idsUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
)

func deleteHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	idsRaw, _ := args["ids"].([]any)
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)

	slog.InfoContext(ctx, "playlistImage delete started")

	var writer bytes.Buffer
	err := del(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlistImage delete failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "playlistImage delete completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func del(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithIDs(ids),
		playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistImage.WithService(nil),
	)

	return pi.Delete(writer)
}

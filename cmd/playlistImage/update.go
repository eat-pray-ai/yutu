package playlistImage

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	updateShort = "Update a playlist image"
	updateLong  = "Update a playlist image for a given playlist id"
)

func init() {
	cmd.MCP.AddTool(updateTool, updateHandler)
	playlistImageCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	updateCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	updateCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	updateCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("playlistId")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var updateTool = mcp.NewTool(
	"playlistImage-update",
	mcp.WithTitleAnnotation(updateShort),
	mcp.WithDescription(updateLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"playlistId", mcp.DefaultString(""),
		mcp.Description(pidUsage), mcp.Required(),
	),
	mcp.WithString(
		"type", mcp.DefaultString(""),
		mcp.Description(typeUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"height", mcp.DefaultNumber(0),
		mcp.Description(heightUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"width", mcp.DefaultNumber(0),
		mcp.Description(widthUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "silent", ""),
		mcp.DefaultString("yaml"), mcp.Description(pkg.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(pkg.JPUsage), mcp.Required(),
	),
)

func updateHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	playlistId, _ = args["playlistId"].(string)
	type_, _ = args["type"].(string)
	heightRaw, _ := args["height"].(float64)
	height = int64(heightRaw)
	widthRaw, _ := args["width"].(float64)
	width = int64(widthRaw)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "playlistImage update started")

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "playlistImage update failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "playlistImage update completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func update(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithPlaylistID(playlistId),
		playlistImage.WithType(type_),
		playlistImage.WithHeight(height),
		playlistImage.WithWidth(width),
		playlistImage.WithMaxResults(1),
		playlistImage.WithService(nil),
	)

	return pi.Update(output, jpath, writer)
}

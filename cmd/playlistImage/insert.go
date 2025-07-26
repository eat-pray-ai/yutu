package playlistImage

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort = "Insert a YouTube playlist image"
	insertLong  = "Insert a YouTube playlist image for a given playlist id"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	playlistImageCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	insertCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	insertCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	insertCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("playlistId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var insertTool = mcp.NewTool(
	"playlistImage-insert",
	mcp.WithTitleAnnotation(insertShort),
	mcp.WithDescription(insertLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"file", mcp.DefaultString(""),
		mcp.Description(fileUsage), mcp.Required(),
	),
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
		mcp.DefaultString("yaml"), mcp.Description(cmd.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JPUsage), mcp.Required(),
	),
)

func insertHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	file, _ = args["file"].(string)
	playlistId, _ = args["playlistId"].(string)
	type_, _ = args["type"].(string)
	heightRaw, _ := args["height"].(float64)
	height = int64(heightRaw)
	widthRaw, _ := args["width"].(float64)
	width = int64(widthRaw)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func insert(writer io.Writer) error {
	pi := playlistImage.NewPlaylistImage(
		playlistImage.WithFile(file),
		playlistImage.WithPlaylistID(playlistId),
		playlistImage.WithType(type_),
		playlistImage.WithHeight(height),
		playlistImage.WithWidth(width),
		playlistImage.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

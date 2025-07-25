package thumbnail

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	cmd.MCP.AddTool(setTool, setHandler)
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	setCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)

	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var setTool = mcp.NewTool(
	"thumbnail-set",
	mcp.WithTitleAnnotation(short),
	mcp.WithDescription(long),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"file", mcp.DefaultString(""),
		mcp.Description(fileUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoId", mcp.DefaultString(""),
		mcp.Description(vidUsage), mcp.Required(),
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

func setHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	file, _ = args["file"].(string)
	videoId, _ = args["videoId"].(string)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	var writer bytes.Buffer
	err := set(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func set(writer io.Writer) error {
	t := thumbnail.NewThumbnail(
		thumbnail.WithFile(file),
		thumbnail.WithVideoId(videoId),
		thumbnail.WithService(nil),
	)

	return t.Set(output, jpath, writer)
}

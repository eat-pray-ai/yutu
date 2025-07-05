package commentThread

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort    = "Insert a new comment thread"
	insertLong     = "Insert a new comment thread"
	insertVidUsage = "ID of the video"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	commentThreadCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringVarP(
		&authorChannelId, "authorChannelId", "a", "", acidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&textOriginal, "textOriginal", "t", "", toUsage)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", insertVidUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JPUsage)

	_ = insertCmd.MarkFlagRequired("authorChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("textOriginal")
	_ = insertCmd.MarkFlagRequired("videoId")
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
	"commentThread-insert",
	mcp.WithTitleAnnotation(insertShort),
	mcp.WithDescription(insertLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"authorChannelId", mcp.DefaultString(""),
		mcp.Description(acidUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"textOriginal", mcp.DefaultString(""),
		mcp.Description(toUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoId", mcp.DefaultString(""),
		mcp.Description(insertVidUsage), mcp.Required(),
	),
	mcp.WithString(
		"output", mcp.Enum("json", "yaml", "silent", ""),
		mcp.DefaultString(""), mcp.Description(cmd.SilentUsage), mcp.Required(),
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
	authorChannelId, _ = args["authorChannelId"].(string)
	channelId, _ = args["channelId"].(string)
	textOriginal, _ = args["textOriginal"].(string)
	videoId, _ = args["videoId"].(string)
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
	ct := commentThread.NewCommentThread(
		commentThread.WithAuthorChannelId(authorChannelId),
		commentThread.WithChannelId(channelId),
		commentThread.WithTextOriginal(textOriginal),
		commentThread.WithVideoId(videoId),
		commentThread.WithService(nil),
	)

	return ct.Insert(output, jpath, writer)
}

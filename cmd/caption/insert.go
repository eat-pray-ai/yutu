package caption

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	insertShort = "Insert caption"
	insertLong  = "Insert caption to a video"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	captionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	insertCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	insertCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	insertCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	insertCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	insertCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	insertCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
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
	"caption-insert",
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
		"audioTrackType",
		mcp.Enum("unknown", "primary", "commentary", "descriptive"),
		mcp.DefaultString("unknown"), mcp.Description(attUsage), mcp.Required(),
	),
	mcp.WithString(
		"isAutoSynced", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(iasUsage), mcp.Required(),
	),
	mcp.WithString(
		"isCC", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(iscUsage), mcp.Required(),
	),
	mcp.WithString(
		"isDraft", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(isdUsage), mcp.Required(),
	),
	mcp.WithString(
		"isEasyReader", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(iserUsage), mcp.Required(),
	),
	mcp.WithString(
		"isLarge", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(islUsage), mcp.Required(),
	),
	mcp.WithString(
		"language", mcp.DefaultString(""),
		mcp.Description(languageUsage), mcp.Required(),
	),
	mcp.WithString(
		"name", mcp.DefaultString(""),
		mcp.Description(nameUsage), mcp.Required(),
	),
	mcp.WithString(
		"trackKind", mcp.Enum("standard", "ASR", "forced"),
		mcp.DefaultString("standard"), mcp.Description(tkUsage), mcp.Required(),
	),
	mcp.WithString(
		"videoId", mcp.DefaultString(""),
		mcp.Description(vidUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOf", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
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
	audioTrackType, _ = args["audioTrackType"].(string)
	isAutoSyncedRaw, _ := args["isAutoSynced"].(string)
	isAutoSynced = utils.BoolPtr(isAutoSyncedRaw)
	isCCRaw, _ := args["isCC"].(string)
	isCC = utils.BoolPtr(isCCRaw)
	isDraftRaw, _ := args["isDraft"].(string)
	isDraft = utils.BoolPtr(isDraftRaw)
	isEasyReaderRaw, _ := args["isEasyReader"].(string)
	isEasyReader = utils.BoolPtr(isEasyReaderRaw)
	isLargeRaw, _ := args["isLarge"].(string)
	isLarge = utils.BoolPtr(isLargeRaw)
	language, _ = args["language"].(string)
	name, _ = args["name"].(string)
	trackKind, _ = args["trackKind"].(string)
	videoId, _ = args["videoId"].(string)
	onBehalfOf, _ = args["onBehalfOf"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "caption insert started")

	var writer bytes.Buffer
	err := insert(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "caption insert failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "caption insert completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func insert(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithFile(file),
		caption.WithAudioTrackType(audioTrackType),
		caption.WithIsAutoSynced(isAutoSynced),
		caption.WithIsCC(isCC),
		caption.WithIsDraft(isDraft),
		caption.WithIsEasyReader(isEasyReader),
		caption.WithIsLarge(isLarge),
		caption.WithLanguage(language),
		caption.WithName(name),
		caption.WithTrackKind(trackKind),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithVideoId(videoId),
		caption.WithService(nil),
	)

	return c.Insert(output, jpath, writer)
}

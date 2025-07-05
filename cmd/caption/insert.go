package caption

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
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
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)

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
		"audioTrackType", mcp.DefaultString("unknown"),
		mcp.Description(attUsage), mcp.Required(),
	),
	mcp.WithBoolean(
		"isAutoSynced", mcp.DefaultBool(true),
		mcp.Description(iasUsage), mcp.Required(),
	),
	mcp.WithBoolean(
		"isCC", mcp.DefaultBool(false),
		mcp.Description(iscUsage), mcp.Required(),
	),
	mcp.WithBoolean(
		"isDraft", mcp.DefaultBool(false),
		mcp.Description(isdUsage), mcp.Required(),
	),
	mcp.WithBoolean(
		"isEasyReader", mcp.DefaultBool(false),
		mcp.Description(iserUsage), mcp.Required(),
	),
	mcp.WithBoolean(
		"isLarge", mcp.DefaultBool(false),
		mcp.Description(islUsage), mcp.Required(),
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
		"trackKind", mcp.DefaultString("standard"),
		mcp.Description(tkUsage), mcp.Required(),
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
		"output", mcp.DefaultString(""),
		mcp.Description(cmd.SilentUsage), mcp.Required(),
	),
	mcp.WithString(
		"jsonpath", mcp.DefaultString(""),
		mcp.Description(cmd.JpUsage), mcp.Required(),
	),
)

func insertHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	file, _ = args["file"].(string)
	audioTrackType, _ = args["audioTrackType"].(string)
	isAutoSyncedRaw, _ := args["isAutoSynced"].(bool)
	isAutoSynced = &isAutoSyncedRaw
	isCCRaw, _ := args["isCC"].(bool)
	isCC = &isCCRaw
	isDraftRaw, _ := args["isDraft"].(bool)
	isDraft = &isDraftRaw
	isEasyReaderRaw, _ := args["isEasyReader"].(bool)
	isEasyReader = &isEasyReaderRaw
	isLargeRaw, _ := args["isLarge"].(bool)
	isLarge = &isLargeRaw
	language, _ = args["language"].(string)
	name, _ = args["name"].(string)
	trackKind, _ = args["trackKind"].(string)
	videoId, _ = args["videoId"].(string)
	onBehalfOf, _ = args["onBehalfOf"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
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

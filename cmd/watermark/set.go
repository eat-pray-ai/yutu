package watermark

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/watermark"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	setShort = "Set watermark for channel's video"
	setLong  = "Set watermark for channel's video by channel id"
)

func init() {
	cmd.MCP.AddTool(setTool, setHandler)
	watermarkCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(
		&inVideoPosition, "inVideoPosition", "p", "", ivpUsage,
	)
	setCmd.Flags().Uint64VarP(&durationMs, "durationMs", "d", 0, dmUsage)
	setCmd.Flags().Uint64VarP(&offsetMs, "offsetMs", "m", 0, omUsage)
	setCmd.Flags().StringVarP(&offsetType, "offsetType", "t", "", otUsage)
	setCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = setCmd.MarkFlagRequired("channelId")
	_ = setCmd.MarkFlagRequired("file")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: setShort,
	Long:  setLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

var setTool = mcp.NewTool(
	"watermark-set",
	mcp.WithTitleAnnotation(setShort),
	mcp.WithDescription(setLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"file", mcp.DefaultString(""),
		mcp.Description(fileUsage), mcp.Required(),
	),
	mcp.WithString(
		"inVideoPosition", mcp.DefaultString(""),
		mcp.Description(ivpUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"durationMs", mcp.DefaultNumber(0),
		mcp.Description(dmUsage), mcp.Required(),
	),
	mcp.WithNumber(
		"offsetMs", mcp.DefaultNumber(0),
		mcp.Description(omUsage), mcp.Required(),
	),
	mcp.WithString(
		"offsetType", mcp.DefaultString(""),
		mcp.Description(otUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
)

func setHandler(
	ctx context.Context, request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	channelId, _ = args["channelId"].(string)
	file, _ = args["file"].(string)
	inVideoPosition, _ = args["inVideoPosition"].(string)
	durationMsRaw, _ := args["durationMs"].(float64)
	durationMs = uint64(durationMsRaw)
	offsetMsRaw, _ := args["offsetMs"].(float64)
	offsetMs = uint64(offsetMsRaw)
	offsetType, _ = args["offsetType"].(string)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)

	var writer bytes.Buffer
	err := set(&writer)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), err
	}
	return mcp.NewToolResultText(writer.String()), nil
}

func set(writer io.Writer) error {
	w := watermark.NewWatermark(
		watermark.WithChannelId(channelId),
		watermark.WithFile(file),
		watermark.WithInVideoPosition(inVideoPosition),
		watermark.WithDurationMs(durationMs),
		watermark.WithOffsetMs(offsetMs),
		watermark.WithOffsetType(offsetType),
		watermark.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		watermark.WithService(nil),
	)

	return w.Set(writer)
}

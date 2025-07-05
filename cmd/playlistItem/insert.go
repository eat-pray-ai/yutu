package playlistItem

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort    = "Insert a playlist item into a playlist"
	insertLong     = "Insert a playlist item into a playlist"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&kind, "kind", "k", "", kindUsage)
	insertCmd.Flags().StringVarP(&kVideoId, "kVideoId", "V", "", kvidUsage)
	insertCmd.Flags().StringVarP(&kChannelId, "kChannelId", "C", "", kcidUsage)
	insertCmd.Flags().StringVarP(&kPlaylistId, "kPlaylistId", "Y", "", kpidUsage)
	insertCmd.Flags().StringVarP(
		&playlistId, "playlistId", "y", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("kind")
	_ = insertCmd.MarkFlagRequired("playlistId")
	_ = insertCmd.MarkFlagRequired("channelId")
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
	"playlistItem-insert",
	mcp.WithTitleAnnotation(insertShort),
	mcp.WithDescription(insertLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"title", mcp.DefaultString(""),
		mcp.Description(titleUsage), mcp.Required(),
	),
	mcp.WithString(
		"description", mcp.DefaultString(""),
		mcp.Description(descUsage), mcp.Required(),
	),
	mcp.WithString(
		"kind", mcp.Enum("video", "channel", "playlist"),
		mcp.DefaultString(""), mcp.Description(kindUsage), mcp.Required(),
	),
	mcp.WithString(
		"kVideoId", mcp.DefaultString(""),
		mcp.Description(kvidUsage), mcp.Required(),
	),
	mcp.WithString(
		"kChannelId", mcp.DefaultString(""),
		mcp.Description(kcidUsage), mcp.Required(),
	),
	mcp.WithString(
		"kPlaylistId", mcp.DefaultString(""),
		mcp.Description(kpidUsage), mcp.Required(),
	),
	mcp.WithString(
		"playlistId", mcp.DefaultString(""),
		mcp.Description(insertPidUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(cidUsage), mcp.Required(),
	),
	mcp.WithString(
		"privacy", mcp.Enum("public", "private", "unlisted"),
		mcp.DefaultString(""), mcp.Description(privacyUsage), mcp.Required(),
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
	title, _ = args["title"].(string)
	description, _ = args["description"].(string)
	kind, _ = args["kind"].(string)
	kVideoId, _ = args["kVideoId"].(string)
	kChannelId, _ = args["kChannelId"].(string)
	kPlaylistId, _ = args["kPlaylistId"].(string)
	playlistId, _ = args["playlistId"].(string)
	channelId, _ = args["channelId"].(string)
	privacy, _ = args["privacy"].(string)
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
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithTitle(title),
		playlistItem.WithDescription(description),
		playlistItem.WithKind(kind),
		playlistItem.WithKVideoId(kVideoId),
		playlistItem.WithKChannelId(kChannelId),
		playlistItem.WithKPlaylistId(kPlaylistId),
		playlistItem.WithPlaylistId(playlistId),
		playlistItem.WithPrivacy(privacy),
		playlistItem.WithChannelId(channelId),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

package video

import (
	"bytes"
	"context"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort     = "Upload a video to YouTube"
	insertLong      = "Upload a video to YouTube, with the specified title, description, tags, etc"
	insertLangUsage = "Language of the video"
)

func init() {
	cmd.MCP.AddTool(insertTool, insertHandler)
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().BoolVarP(
		autoLevels, "autoLevels", "A", true, alUsage,
	)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", insertLangUsage)
	insertCmd.Flags().StringVarP(&license, "license", "L", "youtube", licenseUsage)
	insertCmd.Flags().StringVarP(&thumbnail, "thumbnail", "u", "", thumbnailUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", chidUsage)
	insertCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", pidUsage)
	insertCmd.Flags().StringVarP(&categoryId, "categoryId", "g", "", caidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().BoolVarP(forKids, "forKids", "K", false, fkUsage)
	insertCmd.Flags().BoolVarP(
		embeddable, "embeddable", "E", true, embeddableUsage,
	)
	insertCmd.Flags().StringVarP(&publishAt, "publishAt", "U", "", paUsage)
	insertCmd.Flags().BoolVarP(stabilize, "stabilize", "S", true, stabilizeUsage)
	insertCmd.Flags().BoolVarP(
		notifySubscribers, "notifySubscribers", "N", true, nsUsage,
	)
	insertCmd.Flags().BoolVarP(
		publicStatsViewable, "publicStatsViewable", "P", false, psvUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("categoryId")
	_ = insertCmd.MarkFlagRequired("privacy")
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
	"video-insert",
	mcp.WithTitleAnnotation(insertShort),
	mcp.WithDescription(insertLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithString(
		"autoLevels", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(alUsage), mcp.Required(),
	),
	mcp.WithString(
		"file", mcp.DefaultString(""),
		mcp.Description(fileUsage), mcp.Required(),
	),
	mcp.WithString(
		"title", mcp.DefaultString(""),
		mcp.Description(titleUsage), mcp.Required(),
	),
	mcp.WithString(
		"description", mcp.DefaultString(""),
		mcp.Description(descUsage), mcp.Required(),
	),
	mcp.WithArray(
		"tags", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(tagsUsage), mcp.Required(),
	),
	mcp.WithString(
		"language", mcp.DefaultString(""),
		mcp.Description(insertLangUsage), mcp.Required(),
	),
	mcp.WithString(
		"license", mcp.DefaultString("youtube"),
		mcp.Description(licenseUsage), mcp.Required(),
	),
	mcp.WithString(
		"thumbnail", mcp.DefaultString(""),
		mcp.Description(thumbnailUsage), mcp.Required(),
	),
	mcp.WithString(
		"channelId", mcp.DefaultString(""),
		mcp.Description(chidUsage), mcp.Required(),
	),
	mcp.WithString(
		"playlistId", mcp.DefaultString(""),
		mcp.Description(pidUsage), mcp.Required(),
	),
	mcp.WithString(
		"categoryId", mcp.DefaultString(""),
		mcp.Description(caidUsage), mcp.Required(),
	),
	mcp.WithString(
		"privacy", mcp.DefaultString(""),
		mcp.Description(privacyUsage), mcp.Required(),
	),
	mcp.WithString(
		"forKids", mcp.Enum("true", "false", ""),
		mcp.DefaultString("false"), mcp.Description(fkUsage), mcp.Required(),
	),
	mcp.WithString(
		"embeddable", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(embeddableUsage), mcp.Required(),
	),
	mcp.WithString(
		"publishAt", mcp.DefaultString(""),
		mcp.Description(paUsage), mcp.Required(),
	),
	mcp.WithString(
		"stabilize", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(stabilizeUsage), mcp.Required(),
	),
	mcp.WithString(
		"notifySubscribers", mcp.Enum("true", "false", ""),
		mcp.DefaultString("true"), mcp.Description(nsUsage), mcp.Required(),
	),
	mcp.WithString(
		"publicStatsViewable", mcp.Enum("true", "false", ""),
		mcp.DefaultString("false"), mcp.Description(psvUsage), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwner", mcp.DefaultString(""),
		mcp.Description(""), mcp.Required(),
	),
	mcp.WithString(
		"onBehalfOfContentOwnerChannel", mcp.DefaultString(""),
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
	autoLevelsRaw, _ := args["autoLevels"].(string)
	autoLevels = utils.BoolPtr(autoLevelsRaw)
	file, _ = args["file"].(string)
	title, _ = args["title"].(string)
	description, _ = args["description"].(string)
	tagsRaw, _ := args["tags"].([]any)
	tags = make([]string, len(tagsRaw))
	for i, tag := range tagsRaw {
		tags[i] = tag.(string)
	}
	language, _ = args["language"].(string)
	license, _ = args["license"].(string)
	thumbnail, _ = args["thumbnail"].(string)
	channelId, _ = args["channelId"].(string)
	playListId, _ = args["playlistId"].(string)
	categoryId, _ = args["categoryId"].(string)
	privacy, _ = args["privacy"].(string)
	forKidsRaw, _ := args["forKids"].(string)
	forKids = utils.BoolPtr(forKidsRaw)
	embeddableRaw, _ := args["embeddable"].(string)
	embeddable = utils.BoolPtr(embeddableRaw)
	publishAt, _ = args["publishAt"].(string)
	stabilizeRaw, _ := args["stabilize"].(string)
	stabilize = utils.BoolPtr(stabilizeRaw)
	notifySubscribersRaw, _ := args["notifySubscribers"].(string)
	notifySubscribers = utils.BoolPtr(notifySubscribersRaw)
	publicStatsViewableRaw, _ := args["publicStatsViewable"].(string)
	publicStatsViewable = utils.BoolPtr(publicStatsViewableRaw)
	onBehalfOfContentOwner, _ = args["onBehalfOfContentOwner"].(string)
	onBehalfOfContentOwnerChannel, _ = args["onBehalfOfContentOwnerChannel"].(string)
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
	v := video.NewVideo(
		video.WithAutoLevels(autoLevels),
		video.WithFile(file),
		video.WithTitle(title),
		video.WithDescription(description),
		video.WithTags(tags),
		video.WithLanguage(language),
		video.WithLicense(license),
		video.WithThumbnail(thumbnail),
		video.WithChannelId(channelId),
		video.WithPlaylistId(playListId),
		video.WithCategory(categoryId),
		video.WithPrivacy(privacy),
		video.WithForKids(forKids),
		video.WithEmbeddable(embeddable),
		video.WithPublishAt(publishAt),
		video.WithStabilize(stabilize),
		video.WithNotifySubscribers(notifySubscribers),
		video.WithPublicStatsViewable(publicStatsViewable),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		video.WithService(nil),
	)

	return v.Insert(output, jpath, writer)
}

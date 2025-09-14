package video

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/spf13/cobra"
)

const (
	updateShort     = "Update a video on YouTube"
	updateLong      = "Update a video on YouTube, with the specified title, description, tags, etc"
	updateIdUsage   = "ID of the video to update"
	updateLangUsage = "Language of the video"
)

func init() {
	cmd.MCP.AddTool(updateTool, updateHandler)
	videoCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", updateLangUsage)
	updateCmd.Flags().StringVarP(&license, "license", "L", "youtube", licenseUsage)
	updateCmd.Flags().StringVarP(&thumbnail, "thumbnail", "u", "", thumbnailUsage)
	updateCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", pidUsage)
	updateCmd.Flags().StringVarP(&categoryId, "categoryId", "g", "", caidUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().BoolVarP(
		embeddable, "embeddable", "E", true, embeddableUsage,
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("ids")
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
	"video-update",
	mcp.WithTitleAnnotation(updateShort),
	mcp.WithDescription(updateLong),
	mcp.WithDestructiveHintAnnotation(false),
	mcp.WithOpenWorldHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(false),
	mcp.WithArray(
		"ids", mcp.DefaultArray([]string{}),
		mcp.Items(map[string]any{"type": "string"}),
		mcp.Description(updateIdUsage), mcp.Required(),
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
		mcp.Description(updateLangUsage), mcp.Required(),
	),
	mcp.WithString(
		"license", mcp.Enum("youtube", "creativeCommon"),
		mcp.DefaultString("youtube"), mcp.Description(licenseUsage), mcp.Required(),
	),
	mcp.WithString(
		"thumbnail", mcp.DefaultString(""),
		mcp.Description(thumbnailUsage), mcp.Required(),
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
		"privacy", mcp.Enum("public", "private", "unlisted"),
		mcp.DefaultString(""), mcp.Description(privacyUsage), mcp.Required(),
	),
	mcp.WithString(
		"embeddable", mcp.Enum("true", "false", ""),
		mcp.DefaultString(""), mcp.Description(embeddableUsage), mcp.Required(),
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
	idsRaw, _ := args["ids"].([]any)
	ids = make([]string, len(idsRaw))
	for i, id := range idsRaw {
		ids[i] = id.(string)
	}
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
	playListId, _ = args["playlistId"].(string)
	categoryId, _ = args["categoryId"].(string)
	privacy, _ = args["privacy"].(string)
	embeddableRaw, _ := args["embeddable"].(string)
	embeddable = utils.BoolPtr(embeddableRaw)
	output, _ = args["output"].(string)
	jpath, _ = args["jsonpath"].(string)

	slog.InfoContext(ctx, "video update started")

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "video update failed",
			"error", err,
			"args", args,
		)
		return mcp.NewToolResultError(err.Error()), err
	}
	slog.InfoContext(
		ctx, "video update completed successfully",
		"resultSize", writer.Len(),
	)
	return mcp.NewToolResultText(writer.String()), nil
}

func update(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithTitle(title),
		video.WithDescription(description),
		video.WithTags(tags),
		video.WithLanguage(language),
		video.WithLicense(license),
		video.WithPlaylistId(playListId),
		video.WithThumbnail(thumbnail),
		video.WithCategory(categoryId),
		video.WithPrivacy(privacy),
		video.WithEmbeddable(embeddable),
		video.WithMaxResults(1),
		video.WithService(nil),
	)

	return v.Update(output, jpath, writer)
}

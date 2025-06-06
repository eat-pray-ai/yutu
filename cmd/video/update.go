package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update a video on YouTube"
	updateLong        = "Update a video on YouTube, with the specified title, description, tags, etc"
	updateIdUsage     = "ID of the video to update"
	updateLangUsage   = "Language of the video"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
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
			video.WithService(nil),
		)
		v.Update(output)
	},
}

func init() {
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
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)

	_ = updateCmd.MarkFlagRequired("ids")
}

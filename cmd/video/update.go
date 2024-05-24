package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a video on YouTube",
	Long:  "update a video on YouTube, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithVideoId(id),
			yutuber.WithVideoTitle(title),
			yutuber.WithVideoDescription(description),
			yutuber.WithVideoTags(tags),
			yutuber.WithVideoLanguage(language),
			yutuber.WithVideoLicense(license),
			yutuber.WithVideoPlaylistId(playListId),
			yutuber.WithVideoThumbnail(thumbnail),
			yutuber.WithVideoCategory(categoryId),
			yutuber.WithVideoPrivacy(privacy),
			yutuber.WithVideoEmbeddable(embeddable),
		)
		v.Update()
	},
}

func init() {
	videoCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	updateCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the video")
	updateCmd.Flags().StringArrayVarP(&tags, "tags", "a", []string{}, "Comma separated tags")
	updateCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the video")
	updateCmd.Flags().StringVarP(&license, "license", "L", "youtube", "youtube(default) or creativeCommon")
	updateCmd.Flags().StringVarP(&thumbnail, "thumbnail", "h", "", "Path to the thumbnail")
	updateCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", "Playlist ID of the video")
	updateCmd.Flags().StringVarP(&categoryId, "categoryId", "g", "", "Category of the video")
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "Privacy status of the video: public, private or unlisted")
	updateCmd.Flags().BoolVarP(&embeddable, "embeddable", "e", true, "Whether the video is embeddable")

	updateCmd.MarkFlagRequired("id")
}

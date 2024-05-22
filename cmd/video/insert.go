package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "upload a video to YouTube",
	Long:  "upload a video to YouTube, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithVideoFile(file),
			yutuber.WithVideoTitle(title),
			yutuber.WithVideoDescription(description),
			yutuber.WithVideoTags(tags),
			yutuber.WithVideoLanguage(language),
			yutuber.WithVideoThumbnail(thumbnail),
			yutuber.WithVideoChannelId(channelId),
			yutuber.WithVideoPlaylistId(playListId),
			yutuber.WithVideoCategory(category),
			yutuber.WithVideoPrivacy(privacy),
			yutuber.WithVideoForKids(forKids),
			yutuber.WithVideoEmbeddable(embeddable),
		)
		v.Insert()
	},
}

func init() {
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the video file")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	insertCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the video")
	insertCmd.Flags().StringArrayVarP(&tags, "tags", "a", []string{}, "Comma separated tags")
	insertCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the video")
	insertCmd.Flags().StringVarP(&thumbnail, "thumbnail", "u", "", "Path to the thumbnail")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "Channel ID of the video")
	insertCmd.Flags().StringVarP(&playListId, "playlistId", "y", "", "Playlist ID of the video")
	insertCmd.Flags().StringVarP(&category, "category", "g", "", "Category of the video")
	insertCmd.Flags().StringVarP(
		&privacy, "privacy", "p", "", "Privacy status of the video: public, private, or unlisted",
	)
	insertCmd.Flags().BoolVarP(&forKids, "forKids", "k", false, "Whether the video is for kids")
	insertCmd.Flags().BoolVarP(&embeddable, "embeddable", "e", true, "Whether the video is embeddable")

	insertCmd.MarkFlagRequired("file")
	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("category")
	insertCmd.MarkFlagRequired("privacy")
}

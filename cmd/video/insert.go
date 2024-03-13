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
			yutuber.WithVideoPath(file),
			yutuber.WithVideoTitle(title),
			yutuber.WithVideoDesc(desc),
			yutuber.WithVideoTags(tags),
			yutuber.WithVideoLanguage(language),
			yutuber.WithVideoThumbnail(thumbnail),
			yutuber.WithVideoChannelId(channelId),
			yutuber.WithVideoCategory(category),
			yutuber.WithVideoPrivacy(privacy),
			yutuber.WithVideoForKids(forKids),
			yutuber.WithVideoRestricted(restricted),
			yutuber.WithVideoEmbeddable(embeddable),
		)
		v.Insert()
	},
}

func init() {
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the video file")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	insertCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the video")
	insertCmd.Flags().StringArrayVarP(&tags, "tags", "g", []string{}, "Comma separated tags")
	insertCmd.Flags().StringVar(&language, "language", "", "Language of the video")
	insertCmd.Flags().StringVar(&thumbnail, "thumbnail", "", "Path to the thumbnail")
	insertCmd.Flags().StringVar(&channelId, "channelId", "", "Channel ID of the video")
	insertCmd.Flags().StringVarP(&category, "category", "c", "", "Category of the video")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", "Privacy status of the video")
	insertCmd.Flags().BoolVar(&forKids, "forKids", false, "Whether the video is for kids")
	insertCmd.Flags().BoolVar(&restricted, "restricted", false, "Whether the video is restricted")
	insertCmd.Flags().BoolVar(&embeddable, "embeddable", true, "Whether the video is embeddable")

	insertCmd.MarkFlagRequired("file")
	insertCmd.MarkFlagRequired("title")
	insertCmd.MarkFlagRequired("category")
	insertCmd.MarkFlagRequired("privacy")
}

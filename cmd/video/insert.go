package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "subcommand for inserting a video",
	Long:  `subcommand for inserting a video, which can be used to upload a video to YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithVideoPath(path),
			yutuber.WithVideoTitle(title),
			yutuber.WithVideoDesc(desc),
			yutuber.WithVideoTags(tags),
			yutuber.WithVideoLanguage(language),
			yutuber.WithVideoForKids(forKids),
			yutuber.WithVideoRestricted(restricted),
			yutuber.WithVideoEmbeddable(embeddable),
			yutuber.WithVideoCategory(category),
			yutuber.WithVideoPrivacy(privacy),
			yutuber.WithVideoChannelId(channelId),
		)
		v.Insert()
	},
}

func init() {
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the video file")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	insertCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the video")
	insertCmd.Flags().StringVarP(&tags, "tags", "g", "", "Comma separated tags")
	insertCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the video")
	insertCmd.Flags().BoolVarP(&forKids, "forKids", "f", false, "Whether the video is for kids")
	insertCmd.Flags().BoolVarP(&restricted, "restricted", "e", false, "Whether the video is restricted")
	insertCmd.Flags().BoolVarP(&embeddable, "embeddable", "b", true, "Whether the video is embeddable")
	insertCmd.Flags().StringVarP(&category, "category", "c", "", "Category of the video")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "r", "", "Privacy status of the video")
	insertCmd.Flags().StringVarP(&channelId, "channelId", "i", "", "Channel ID of the video")
}

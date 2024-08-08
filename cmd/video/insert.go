package video

import (
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Upload a video to YouTube",
	Long:  "Upload a video to YouTube, with the specified title, description, tags, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		v := video.NewVideo(
			video.WithAutoLevels(autoLevels, true),
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
			video.WithStabilize(stabilize, true),
			video.WithNotifySubscribers(notifySubscribers),
			video.WithPublicStatsViewable(publicStatsViewable),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			video.WithService(nil),
		)
		v.Insert(output)
	},
}

func init() {
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().BoolVarP(
		&autoLevels, "autoLevels", "a", true, "Should auto-levels be applied to the upload",
	)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the video file")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	insertCmd.Flags().StringVarP(
		&description, "description", "d", "", "Description of the video",
	)
	insertCmd.Flags().StringArrayVarP(
		&tags, "tags", "T", []string{}, "Comma separated tags",
	)
	insertCmd.Flags().StringVarP(
		&language, "language", "l", "", "Language of the video",
	)
	insertCmd.Flags().StringVarP(
		&license, "license", "L", "youtube", "youtube(default) or creativeCommon",
	)
	insertCmd.Flags().StringVarP(
		&thumbnail, "thumbnail", "u", "", "Path to the thumbnail",
	)
	insertCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", "Channel ID of the video",
	)
	insertCmd.Flags().StringVarP(
		&playListId, "playlistId", "y", "", "Playlist ID of the video",
	)
	insertCmd.Flags().StringVarP(
		&categoryId, "categoryId", "g", "", "Category of the video",
	)
	insertCmd.Flags().StringVarP(
		&privacy, "privacy", "p", "",
		"Privacy status of the video: public, private, or unlisted",
	)
	insertCmd.Flags().BoolVarP(
		&forKids, "forKids", "k", false, "Whether the video is for kids",
	)
	insertCmd.Flags().BoolVarP(
		&embeddable, "embeddable", "e", true, "Whether the video is embeddable",
	)
	insertCmd.Flags().StringVarP(
		&publishAt, "publishAt", "U", "",
		"Datetime when the video is scheduled to publish",
	)
	insertCmd.Flags().BoolVarP(&stabilize, "stabilize", "s", true, "Should stabilize be applied to the upload")
	insertCmd.Flags().BoolVarP(
		&notifySubscribers, "notifySubscribers", "n", true,
		"Notify the channel subscribers about the new video",
	)
	insertCmd.Flags().BoolVarP(
		&publicStatsViewable, "publicStatsViewable", "P", false,
		"Whether the extended video statistics can be viewed by everyone",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")

	insertCmd.MarkFlagRequired("file")
	insertCmd.MarkFlagRequired("categoryId")
	insertCmd.MarkFlagRequired("privacy")
}

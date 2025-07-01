package video

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/spf13/cobra"
)

const (
	insertShort     = "Upload a video to YouTube"
	insertLong      = "Upload a video to YouTube, with the specified title, description, tags, etc"
	insertLangUsage = "Language of the video"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
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

		err := v.Insert(output, jpath, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
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

package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update caption",
	Long:  "Update caption of a video",
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithFile(file),
			caption.WithAudioTrackType(audioTrackType),
			caption.WithIsAutoSynced(isAutoSynced, cmd.Flags().Lookup("isAutoSynced").Changed),
			caption.WithIsCC(isCC, cmd.Flags().Lookup("isCC").Changed),
			caption.WithIsDraft(isDraft, cmd.Flags().Lookup("isDraft").Changed),
			caption.WithIsEasyReader(isEasyReader, cmd.Flags().Lookup("isEasyReader").Changed),
			caption.WithIsLarge(isLarge, cmd.Flags().Lookup("isLarge").Changed),
			caption.WithLanguage(language),
			caption.WithName(name),
			caption.WithTrackKind(trackKind),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithVideoId(videoId),
			caption.WithService(nil),
		)
		c.Update(output)
	},
}

func init() {
	captionCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the caption file")
	updateCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", "unknown, primary, commentary or descriptive",
	)
	updateCmd.Flags().BoolVarP(
		&isAutoSynced, "isAutoSynced", "A", true,
		"Whether YouTube synchronized the caption track to the audio track in the video",
	)
	updateCmd.Flags().BoolVarP(
		&isCC, "isCC", "C", false,
		"Whether the track contains closed captions for the deaf and hard of hearing",
	)
	updateCmd.Flags().BoolVarP(&isDraft, "isDraft", "D", false, "whether the caption track is a draft")
	updateCmd.Flags().BoolVarP(
		&isEasyReader, "isEasyReader", "E", false, "Whether caption track is formatted for 'easy reader'",
	)
	updateCmd.Flags().BoolVarP(
		&isLarge, "isLarge", "L", false,
		"Whether the caption track uses large text for the vision-impaired",
	)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the caption track")
	updateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the caption track")
	updateCmd.Flags().StringVarP(&trackKind, "trackKind", "t", "standard", "standard, ASR or forced")
	updateCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	updateCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	updateCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "")
	updateCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

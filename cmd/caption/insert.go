package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert caption",
	Long:  "Insert caption to a video",
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithFile(file),
			caption.WithAudioTrackType(audioTrackType),
			caption.WithIsAutoSynced(isAutoSynced, true),
			caption.WithIsCC(isCC, true),
			caption.WithIsDraft(isDraft, true),
			caption.WithIsEasyReader(isEasyReader, true),
			caption.WithIsLarge(isLarge, true),
			caption.WithLanguage(language),
			caption.WithName(name),
			caption.WithTrackKind(trackKind),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithVideoId(videoId),
			caption.WithService(nil),
		)
		c.Insert(output)
	},
}

func init() {
	captionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the caption file")
	insertCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", "unknown, primary, commentary or descriptive",
	)
	insertCmd.Flags().BoolVarP(
		&isAutoSynced, "isAutoSynced", "A", true,
		"Whether YouTube synchronized the caption track to the audio track in the video",
	)
	insertCmd.Flags().BoolVarP(
		&isCC, "isCC", "C", false,
		"Whether the track contains closed captions for the deaf and hard of hearing",
	)
	insertCmd.Flags().BoolVarP(&isDraft, "isDraft", "D", false, "whether the caption track is a draft")
	insertCmd.Flags().BoolVarP(
		&isEasyReader, "isEasyReader", "E", false, "Whether caption track is formatted for 'easy reader'",
	)
	insertCmd.Flags().BoolVarP(
		&isLarge, "isLarge", "L", false,
		"Whether the caption track uses large text for the vision-impaired",
	)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", "Language of the caption track")
	insertCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the caption track")
	insertCmd.Flags().StringVarP(&trackKind, "trackKind", "t", "standard", "standard, ASR or forced")
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	insertCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	insertCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "")
	insertCmd.Flags().StringVarP(&output, "output", "o", "", "json, yaml or silent")
}

package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update caption"
	updateLong        = "Update caption of a video"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithFile(file),
			caption.WithAudioTrackType(audioTrackType),
			caption.WithIsAutoSynced(isAutoSynced),
			caption.WithIsCC(isCC),
			caption.WithIsDraft(isDraft),
			caption.WithIsEasyReader(isEasyReader),
			caption.WithIsLarge(isLarge),
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

	updateCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	updateCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	updateCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	updateCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	updateCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	updateCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	updateCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	updateCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	updateCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	updateCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)
}

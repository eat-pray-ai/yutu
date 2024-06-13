package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/caption"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list captions",
	Long:  "list captions of a video",
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithId(id),
			caption.WithVideoId(videoId),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		c.List(parts, output)
	},
}

func init() {
	captionCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the caption")
	listCmd.Flags().StringVarP(&videoId, "videoId", "I", "", "ID of the video")
	listCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json or yaml")
}

package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List captions",
	Long:  "List captions of a video",
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithID(id),
			caption.WithVideoId(videoId),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithService(nil),
		)
		c.List(parts, output)
	},
}

func init() {
	captionCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the caption")
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", "ID of the video")
	listCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json or yaml")
}

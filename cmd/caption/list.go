package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

const (
	listOutputUsage = "json or yaml"
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

	listCmd.Flags().StringVarP(&id, "id", "i", "", idUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
}

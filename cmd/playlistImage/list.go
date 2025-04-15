package playlistImage

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List YouTube playlist images",
	Long:  "List YouTube playlist images' info",
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithParent(parent),
			playlistImage.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			playlistImage.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			playlistImage.WithMaxResults(maxResults),
			playlistImage.WithService(nil),
		)
		pi.List(parts, output)
	},
}

func init() {
	playlistImageCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&parent, "parent", "p", "", "Return PlaylistImages for this playlist id")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	listCmd.Flags().StringVarP(&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "")
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned")
	listCmd.Flags().StringSliceVarP(&parts, "parts", "p", []string{"id", "kind", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json or yaml")
}

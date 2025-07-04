package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List playlist's info"
	listLong     = "List playlist's info, such as title, description, etc"
	listIdsUsage = "Return the playlists with the given IDs for Stubby or Apiary"
	listCidUsage = "Return the playlists owned by the specified channel id"
)

func init() {
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func list(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithIDs(ids),
		playlist.WithChannelId(channelId),
		playlist.WithHl(hl),
		playlist.WithMaxResults(maxResults),
		playlist.WithMine(mine),
		playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlist.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		playlist.WithService(nil),
	)

	return p.List(parts, output, jpath, writer)
}

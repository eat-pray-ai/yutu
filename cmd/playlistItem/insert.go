package playlistItem

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort    = "Insert a playlist item into a playlist"
	insertLong     = "Insert a playlist item into a playlist"
	insertPidUsage = "The id that YouTube uses to uniquely identify the playlist that the item is in"
)

func init() {
	playlistItemCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&kind, "kind", "k", "", kindUsage)
	insertCmd.Flags().StringVarP(&kVideoId, "kVideoId", "V", "", kvidUsage)
	insertCmd.Flags().StringVarP(&kChannelId, "kChannelId", "C", "", kcidUsage)
	insertCmd.Flags().StringVarP(&kPlaylistId, "kPlaylistId", "Y", "", kpidUsage)
	insertCmd.Flags().StringVarP(
		&playlistId, "playlistId", "y", "", insertPidUsage,
	)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("kind")
	_ = insertCmd.MarkFlagRequired("playlistId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insert(writer io.Writer) error {
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithTitle(title),
		playlistItem.WithDescription(description),
		playlistItem.WithKind(kind),
		playlistItem.WithKVideoId(kVideoId),
		playlistItem.WithKChannelId(kChannelId),
		playlistItem.WithKPlaylistId(kPlaylistId),
		playlistItem.WithPlaylistId(playlistId),
		playlistItem.WithPrivacy(privacy),
		playlistItem.WithChannelId(channelId),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.Insert(output, jpath, writer)
}

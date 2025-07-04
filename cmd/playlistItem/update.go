package playlistItem

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/spf13/cobra"
	"io"
)

const (
	updateShort   = "Update a playlist item"
	updateLong    = "Update a playlist item's info, such as title, description, etc"
	updateIdUsage = "ID of the playlist item to update"
)

func init() {
	playlistItemCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func update(writer io.Writer) error {
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithIDs(ids),
		playlistItem.WithTitle(title),
		playlistItem.WithDescription(description),
		playlistItem.WithPrivacy(privacy),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithService(nil),
	)

	return pi.Update(output, jpath, writer)
}

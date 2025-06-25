package playlistImage

import (
	"github.com/eat-pray-ai/yutu/pkg/playlistImage"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update a playlist image"
	updateLong        = "Update a playlist image for a given playlist id"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		pi := playlistImage.NewPlaylistImage(
			playlistImage.WithPlaylistID(playlistId),
			playlistImage.WithType(type_),
			playlistImage.WithHeight(height),
			playlistImage.WithWidth(width),
			playlistImage.WithService(nil),
		)

		err := pi.Update(output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	playlistImageCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&playlistId, "playlistId", "p", "", pidUsage)
	updateCmd.Flags().StringVarP(&type_, "type", "t", "", typeUsage)
	updateCmd.Flags().Int64VarP(&height, "height", "H", 0, heightUsage)
	updateCmd.Flags().Int64VarP(&width, "width", "W", 0, widthUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)

	_ = updateCmd.MarkFlagRequired("playlistId")
}

package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list video's info",
	Long:  "list video's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithVideoId(id),
		)
		v.List()
	},
}

func init() {
	videoCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video")
	listCmd.MarkFlagRequired("id")
}

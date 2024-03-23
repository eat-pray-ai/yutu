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
			yutuber.WithVideoRating(rating),
		)
		v.List(parts, output)
	},
}

func init() {
	videoCmd.AddCommand(listCmd)
	parts := []string{"id", "snippet", "status"}

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the video")
	listCmd.Flags().StringVarP(&rating, "rating", "r", "", "My rating of the video: like or dislike")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", parts, "Comma separated parts")
}

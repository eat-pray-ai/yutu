package video

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "subcommand for inserting a video",
	Long:  `subcommand for inserting a video, which can be used to upload a video to YouTube.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := yutuber.NewVideo(
			yutuber.WithPath(path),
			yutuber.WithTitle(title),
			yutuber.WithDesc(desc),
			yutuber.WithCategory(category),
			yutuber.WithKeywords(keywords),
			yutuber.WithPrivacy(privacy),
		)
		v.Insert()
	},
}

func init() {
	videoCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&path, "path", "p", "", "Path to the video file")
	insertCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the video")
	insertCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the video")
	insertCmd.Flags().StringVarP(&category, "category", "c", "", "Category of the video")
	insertCmd.Flags().StringVarP(&keywords, "keywords", "k", "", "Comma separated keywords")
	insertCmd.Flags().StringVarP(&privacy, "privacy", "r", "", "Privacy status of the video")
}

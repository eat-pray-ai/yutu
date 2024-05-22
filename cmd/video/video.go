package video

import (
	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

var (
	id          string
	file        string
	title       string
	description string
	tags        []string
	language    string
	thumbnail   string
	rating      string
	chart       string
	channelId   string
	playListId  string
	category    string
	privacy     string
	forKids     bool
	embeddable  bool
	output      string
	parts       []string
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "manipulate YouTube videos",
	Long:  "manipulate YouTube videos, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// videoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// videoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

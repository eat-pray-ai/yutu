package comment

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	IDs              []string
	AuthorChannelId  string
	CanRate          bool
	ChannelId        string
	MaxResults       int64
	ParentId         string
	TextFormat       string
	TextOriginal     string
	ModerationStatus string
	BanAuthor        bool
	VideoId          string
	ViewerRating     string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manipulate YouTube comments",
	Long:  "Manipulate YouTube comments, such as insert, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(commentCmd)
}

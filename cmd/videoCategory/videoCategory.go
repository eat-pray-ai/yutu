package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id         string
	hl         string
	regionCode string
	parts      []string
	output     string
)

var videoCategoryCmd = &cobra.Command{
	Use:   "videoCategory",
	Short: "list YouTube video categories",
	Long:  "list YouTube video categories",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCategoryCmd)
}

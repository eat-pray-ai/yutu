package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	regionCode string
	parts      []string
	output     string
)

var videoCategoryCmd = &cobra.Command{
	Use:   "videoCategory",
	Short: "manipulate YouTube video categories",
	Long:  "manipulate YouTube video categories, only list for now",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCategoryCmd)
}

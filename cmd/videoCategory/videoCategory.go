package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short      = "List YouTube video categories"
	long       = "List YouTube video categories' info, such as id, title, assignable, etc."
	idsUsage   = "IDs of the video categories"
	hlUsage    = "Host language"
	rcUsage    = "Region code"
	partsUsage = "Comma separated parts"
)

var (
	ids        []string
	hl         string
	regionCode string
	parts      []string
	output     string
	jpath      string
)

var videoCategoryCmd = &cobra.Command{
	Use:   "videoCategory",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCategoryCmd)
}

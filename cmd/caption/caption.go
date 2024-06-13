package caption

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id                     string
	file                   string
	audioTrackType         string
	isAutoSynced           *bool
	isCC                   *bool
	isDraft                *bool
	isEasyReader           *bool
	isLarge                *bool
	language               string
	name                   string
	trackKind              string
	onBehalfOf             string
	onBehalfOfContentOwner string
	videoId                string
	parts                  []string
	output                 string
)

var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: "manipulate YouTube captions",
	Long:  "manipulate YouTube captions, such as list, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(captionCmd)
}

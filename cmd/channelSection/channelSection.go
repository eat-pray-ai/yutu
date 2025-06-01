package channelSection

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	id                     string
	channelId              string
	hl                     string
	mine                   = utils.BoolPtr("false")
	onBehalfOfContentOwner string
	parts                  []string
	output                 string
)

var channelSectionCmd = &cobra.Command{
	Use:   "channelSection",
	Short: "Manipulate YouTube channel sections",
	Long:  "List or delete YouTube channel sections",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.ResetBool(map[string]*bool{"mine": mine}, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelSectionCmd)
}

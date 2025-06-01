package caption

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	id                     string
	file                   string
	audioTrackType         string
	isAutoSynced           = utils.BoolPtr("false")
	isCC                   = utils.BoolPtr("false")
	isDraft                = utils.BoolPtr("false")
	isEasyReader           = utils.BoolPtr("false")
	isLarge                = utils.BoolPtr("false")
	language               string
	name                   string
	trackKind              string
	onBehalfOf             string
	onBehalfOfContentOwner string
	videoId                string
	parts                  []string
	output                 string
	tfmt                   string
	tlang                  string
)

var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: "Manipulate YouTube captions",
	Long:  "List, insert, update, download, or delete YouTube captions",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resetFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(captionCmd)
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]*bool{
		"isAutoSynced": isAutoSynced,
		"isCC":         isCC,
		"isDraft":      isDraft,
		"isEasyReader": isEasyReader,
		"isLarge":      isLarge,
	}

	utils.ResetBool(boolMap, flagSet)
}

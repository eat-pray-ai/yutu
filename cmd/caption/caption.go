package caption

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	short         = "Manipulate YouTube captions"
	long          = "List, insert, update, download, or delete YouTube captions"
	fileUsage     = "Path to save the caption file"
	attUsage      = "unknown, primary, commentary, or descriptive"
	iasUsage      = "Whether YouTube synchronized the caption track to the audio track in the video"
	iscUsage      = "Whether the track contains closed captions for the deaf and hard of hearing"
	isdUsage      = "whether the caption track is a draft"
	iserUsage     = "Whether caption track is formatted for 'easy reader'"
	islUsage      = "Whether the caption track uses large text for the vision-impaired"
	languageUsage = "Language of the caption track"
	nameUsage     = "Name of the caption track"
	tkUsage       = "standard, ASR, or forced"
	vidUsage      = "ID of the video"
	tfmtUsage     = "sbv, srt, or vtt"
	tlangUsage    = "Translate the captions into this language"
)

var (
	ids                    []string
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
	tfmt                   string
	tlang                  string
	output                 string
	jpath                  string
)

var captionCmd = &cobra.Command{
	Use:   "caption",
	Short: short,
	Long:  long,
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
	boolMap := map[string]**bool{
		"isAutoSynced": &isAutoSynced,
		"isCC":         &isCC,
		"isDraft":      &isDraft,
		"isEasyReader": &isEasyReader,
		"isLarge":      &isLarge,
	}

	utils.ResetBool(boolMap, flagSet)
}

package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

var (
	categoryId             string
	forHandle              string
	forUsername            string
	hl                     string
	id                     string
	managedByMe            = utils.BoolPtr("")
	maxResults             int64
	mine                   = utils.BoolPtr("")
	mySubscribers          = utils.BoolPtr("")
	onBehalfOfContentOwner string

	country         string
	customUrl       string
	defaultLanguage string
	description     string
	title           string
	output          string
	parts           []string
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Manipulate YouTube channels",
	Long:  "List or update YouTube channels",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resetFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// channelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// channelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]*bool{
		"managedByMe":   managedByMe,
		"mine":          mine,
		"mySubscribers": mySubscribers,
	}

	utils.ResetBool(boolMap, flagSet)
}

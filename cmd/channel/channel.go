package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"

	"github.com/spf13/cobra"
)

var (
	categoryId             string
	forHandle              string
	forUsername            string
	hl                     string
	id                     string
	managedByMe            string
	maxResults             int64
	mine                   string
	mySubscribers          string
	onBehalfOfContentOwner string

	title  string
	desc   string
	output string
	parts  []string
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "manipulate YouTube channels",
	Long:  "manipulate YouTube channels, such as list, update, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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

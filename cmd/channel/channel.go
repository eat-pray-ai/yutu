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
	managedByMe            bool
	maxResults             int64
	mine                   bool
	mySubscribers          bool
	onBehalfOfContentOwner string

	country         string
	customUrl       string
	defaultLanguage string
	description     string
	title           string
	output          string
	parts           []string
	credential      string
	cacheToken      string
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Manipulate YouTube channels",
	Long:  "List or update YouTube channels",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelCmd)

	channelCmd.PersistentFlags().StringVarP(&credential, "credential", "c", "client_secret.json", "Path to client secret file")
	channelCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file")
}

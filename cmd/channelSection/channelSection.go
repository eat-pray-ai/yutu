package channelSection

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	id                     string
	channelId              string
	hl                     string
	mine                   bool
	onBehalfOfContentOwner string
	parts                  []string
	output                 string
	credential             string
	cacheToken             string
)

var channelSectionCmd = &cobra.Command{
	Use:   "channelSection",
	Short: "Manipulate channel section",
	Long:  "List or delete channel section",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(channelSectionCmd)

	channelSectionCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	channelSectionCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

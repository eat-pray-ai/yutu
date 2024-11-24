package watermark

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var (
	channelId              string
	file                   string
	inVideoPosition        string
	durationMs             uint64
	offsetMs               uint64
	offsetType             string
	onBehalfOfContentOwner string
	credential             string
	cacheToken             string
)

var wartermarkCmd = &cobra.Command{
	Use:   "watermark",
	Short: "Manipulate Youtube watermarks",
	Long:  "Set or unset Youtube watermarks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(wartermarkCmd)

	wartermarkCmd.PersistentFlags().StringVarP(&credential, "credential", "", "client_secret.json", "Path to client secret file")
	wartermarkCmd.PersistentFlags().StringVarP(&cacheToken, "cacheToken", "", "youtube.token.json", "Path to token cache file")
}

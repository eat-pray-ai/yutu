package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list channel's info",
	Long:  `list channel's info, such as title, description, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := yutuber.NewChannel(
			yutuber.WithChannelID(id),
		)
		c.List()
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the channel")
	listCmd.MarkFlagRequired("id")
}

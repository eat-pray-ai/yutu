package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update channel's info",
	Long:  "update channel's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		c := yutuber.NewChannel(
			yutuber.WithChannelID(id),
			yutuber.WithChannelTitle(title),
			yutuber.WithChannelDesc(desc),
		)
		c.Update()
	},
}

func init() {
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the channel")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the channel")
	updateCmd.Flags().StringVarP(&desc, "desc", "d", "", "Description of the channel")

	updateCmd.MarkFlagRequired("id")
}

package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list channel's info",
	Long:  "list channel's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		c := yutuber.NewChannel(
			yutuber.WithChannelID(id),
			yutuber.WithChannelUser(user),
		)
		c.List(parts, output)
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&id, "id", "i", "", "Return the channels with the specified IDs")
	listCmd.Flags().StringVarP(&user, "user", "u", "", "Return the channel associated with a YouTube username")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet", "status"}, "Comma separated parts")
}

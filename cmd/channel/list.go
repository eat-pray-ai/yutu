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
			yutuber.WithChannelCategoryId(categoryId),
			yutuber.WithChannelForHandle(forHandle),
			yutuber.WithChannelForUsername(forUsername),
			yutuber.WithChannelHl(hl),
			yutuber.WithChannelId(id),
			yutuber.WithChannelManagedByMe(managedByMe),
			yutuber.WithChannelMaxResults(maxResults),
			yutuber.WithChannelMine(mine),
			yutuber.WithChannelMySubscribers(mySubscribers),
			yutuber.WithChannelOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		c.List(parts, output)
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&categoryId, "categoryId", "c", "", "Return the channels within the specified guide category ID",
	)
	listCmd.Flags().StringVarP(&forHandle, "forHandle", "d", "", "Return the channel associated with a YouTube handle")
	listCmd.Flags().StringVarP(
		&forUsername, "forUsername", "u", "", "Return the channel associated with a YouTube username",
	)
	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Specifies the localization language of the metadata",
	)
	listCmd.Flags().StringVarP(&id, "id", "i", "", "Return the channels with the specified IDs")
	listCmd.Flags().StringVarP(
		&managedByMe, "managedByMe", "M", "", "Specify the maximum number of items that should be returned",
	)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned")
	listCmd.Flags().StringVarP(&mine, "mine", "m", "", "Return the ids of channels owned by the authenticated user")
	listCmd.Flags().StringVarP(
		&mySubscribers, "mySubscribers", "s", "", "Return the channels subscribed to the authenticated user",
	)
	listCmd.Flags().StringVarP(&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet", "status"}, "Comma separated parts")
}

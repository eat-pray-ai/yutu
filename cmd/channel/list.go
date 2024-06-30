package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/channel"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List channel's info",
	Long:  "List channel's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		c := channel.NewChannel(
			channel.WithCategoryId(categoryId),
			channel.WithForHandle(forHandle),
			channel.WithForUsername(forUsername),
			channel.WithHl(hl),
			channel.WithID(id),
			channel.WithChannelManagedByMe(managedByMe, true),
			channel.WithMaxResults(maxResults),
			channel.WithMine(mine, true),
			channel.WithMySubscribers(mySubscribers, true),
			channel.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channel.WithService(nil),
		)
		c.List(parts, output)
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&categoryId, "categoryId", "c", "",
		"Return the channels within the specified guide category ID",
	)
	listCmd.Flags().StringVarP(
		&forHandle, "forHandle", "d", "",
		"Return the channel associated with a YouTube handle",
	)
	listCmd.Flags().StringVarP(
		&forUsername, "forUsername", "u", "",
		"Return the channel associated with a YouTube username",
	)
	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Specifies the localization language of the metadata",
	)
	listCmd.Flags().StringVarP(
		&id, "id", "i", "", "Return the channels with the specified IDs",
	)
	listCmd.Flags().BoolVarP(
		&managedByMe, "managedByMe", "M", false,
		"Specify the maximum number of items that should be returned",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5,
		"The maximum number of items that should be returned",
	)
	listCmd.Flags().BoolVarP(
		&mine, "mine", "m", true,
		"Return the ids of channels owned by the authenticated user",
	)
	listCmd.Flags().BoolVarP(
		&mySubscribers, "mySubscribers", "s", false,
		"Return the channels subscribed to the authenticated user",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"},
		"Comma separated parts",
	)
}

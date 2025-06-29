package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/spf13/cobra"
)

const (
	listShort       = "List channel's info"
	listLong        = "List channel's info, such as title, description, etc."
	listIdsUsage    = "Return the channels with the specified IDs"
	listOutputUsage = "json or yaml"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := channel.NewChannel(
			channel.WithCategoryId(categoryId),
			channel.WithForHandle(forHandle),
			channel.WithForUsername(forUsername),
			channel.WithHl(hl),
			channel.WithIDs(ids),
			channel.WithChannelManagedByMe(managedByMe),
			channel.WithMaxResults(maxResults),
			channel.WithMine(mine),
			channel.WithMySubscribers(mySubscribers),
			channel.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channel.WithService(nil),
		)

		err := c.List(parts, output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	channelCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&categoryId, "categoryId", "g", "", cidUsage,
	)
	listCmd.Flags().StringVarP(
		&forHandle, "forHandle", "d", "", fhUsage,
	)
	listCmd.Flags().StringVarP(
		&forUsername, "forUsername", "u", "", fuUsage,
	)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().BoolVarP(
		managedByMe, "managedByMe", "E", false, mbmUsage,
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, mrUsage,
	)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().BoolVarP(
		mySubscribers, "mySubscribers", "S", false, msUsage,
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", listOutputUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
}

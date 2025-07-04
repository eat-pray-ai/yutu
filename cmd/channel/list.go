package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/spf13/cobra"
	"io"
)

const (
	listShort    = "List channel's info"
	listLong     = "List channel's info, such as title, description, etc."
	listIdsUsage = "Return the channels with the specified IDs"
)

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
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", cmd.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func list(writer io.Writer) error {
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

	return c.List(parts, output, jpath, writer)
}

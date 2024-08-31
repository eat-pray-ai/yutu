package superChatEvent

import (
	"github.com/eat-pray-ai/yutu/pkg/superChatEvent"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Super Chat events for a channel",
	Long:  "List Super Chat events for a channel",
	Run: func(cmd *cobra.Command, args []string) {
		sc := superChatEvent.NewSuperChatEvent(
			superChatEvent.WithHl(hl),
			superChatEvent.WithMaxResults(maxResults),
		)
		sc.List(parts, output)
	},
}

func init() {
	superChatEventCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&hl, "hl", "l", "", "Return rendered funding amounts in specified language",
	)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, "The maximum number of items that should be returned",
	)
	listCmd.Flags().StringSliceVarP(&parts, "parts", "p", []string{"id", "snippet"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "json or yaml")
}

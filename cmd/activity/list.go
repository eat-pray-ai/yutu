package activity

import (
	"github.com/eat-pray-ai/yutu/pkg/activity"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List YouTube activities",
	Long:  "List YouTube activities, such as likes, favorites, uploads, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		a := activity.NewActivity(
			activity.WithChannelId(channelId),
			activity.WithHome(home, true),
			activity.WithMaxResults(maxResults),
			activity.WithMine(mine, true),
			activity.WithPublishedAfter(publishedAfter),
			activity.WithPublishedBefore(publishedBefore),
			activity.WithRegionCode(regionCode),
			activity.WithService(nil),
		)
		a.List(parts, output)
	},
}

func init() {
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(
		&channelId, "channelId", "c", "", "ID of the channel",
	)
	listCmd.Flags().BoolVarP(&home, "home", "H", true, "true or false")
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "x", 5,
		"Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().BoolVarP(&mine, "mine", "m", true, "true or false")
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "",
		"Filter on activities published after this date",
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "",
		"Filter on activities published before this date",
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", "")

	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "contentDetails"},
		"Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

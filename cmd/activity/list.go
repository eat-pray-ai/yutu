package activity

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list activities",
	Long:  "list activities, such as likes, favorites, uploads, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		a := yutuber.NewActivity(
			yutuber.WithActivityChannelId(channelId),
			yutuber.WithActivityHome(home),
			yutuber.WithActivityMaxResults(maxResults),
			yutuber.WithActivityMine(mine),
			yutuber.WithActivityPublishedAfter(publishedAfter),
			yutuber.WithActivityPublishedBefore(publishedBefore),
			yutuber.WithActivityRegionCode(regionCode),
		)
		a.List(parts, output)
	},
}

func init() {
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of the channel")
	listCmd.Flags().StringVarP(&home, "home", "h", "", "true or false")
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "x", 5, "Specifies the maximum number of items that should be returned",
	)
	listCmd.Flags().StringVarP(&mine, "mine", "m", "", "true or false")
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", "Filter on activities published after this date",
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", "Filter on activities published before this date",
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", "")

	listCmd.Flags().StringArrayVarP(
		&parts, "parts", "p", []string{"id", "snippet", "contentDetails"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
}

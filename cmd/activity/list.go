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
		)
		a.List(parts, output)
	},
}

func init() {
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", "ID of the channel")
	listCmd.Flags().StringArrayVarP(&parts, "parts", "p", []string{"id", "snippet", "contentDetails"}, "Comma separated parts")
	listCmd.Flags().StringVarP(&output, "output", "o", "", "Output format: json or yaml")
}

package subscription

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/subscription"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort    = "Insert a YouTube subscription"
	insertLong     = "Insert a YouTube subscription"
	insertCidUsage = "ID of the channel to be subscribed"
)

func init() {
	subscriptionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(
		&subscriberChannelId, "subscriberChannelId", "s", "", scidUsage,
	)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("subscriberChannelId")
	_ = insertCmd.MarkFlagRequired("channelId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insert(writer io.Writer) error {
	s := subscription.NewSubscription(
		subscription.WithSubscriberChannelId(subscriberChannelId),
		subscription.WithDescription(description),
		subscription.WithChannelId(channelId),
		subscription.WithTitle(title),
		subscription.WithService(nil),
	)

	return s.Insert(output, jpath, writer)
}

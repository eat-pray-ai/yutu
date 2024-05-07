package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update channel's info",
	Long:  "update channel's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		c := yutuber.NewChannel(
			yutuber.WithChannelId(id),
			yutuber.WithChannelCountry(country),
			yutuber.WithChannelCustomUrl(customUrl),
			yutuber.WithChannelDefaultLanguage(defaultLanguage),
			yutuber.WithChannelDescription(description),
			yutuber.WithChannelTitle(title),
		)
		c.Update()
	},
}

func init() {
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the channel")
	updateCmd.Flags().StringVarP(&country, "country", "c", "", "Country of the channel")
	updateCmd.Flags().StringVarP(&customUrl, "customUrl", "u", "", "Custom URL of the channel")
	updateCmd.Flags().StringVarP(
		&defaultLanguage, "defaultLanguage", "l", "", "The language of the channel's default title and description",
	)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the channel")
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the channel")

	updateCmd.MarkFlagRequired("id")
}

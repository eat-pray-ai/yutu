package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/channel"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update channel's info",
	Long:  "Update channel's info, such as title, description, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		c := channel.NewChannel(
			channel.WithID(id),
			channel.WithCountry(country),
			channel.WithCustomUrl(customUrl),
			channel.WithDefaultLanguage(defaultLanguage),
			channel.WithDescription(description),
			channel.WithTitle(title),
			channel.WithService(nil),
		)
		c.Update()
	},
}

func init() {
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the channel")
	updateCmd.Flags().StringVarP(
		&country, "country", "c", "", "Country of the channel",
	)
	updateCmd.Flags().StringVarP(
		&customUrl, "customUrl", "u", "", "Custom URL of the channel",
	)
	updateCmd.Flags().StringVarP(
		&defaultLanguage, "defaultLanguage", "l", "",
		"The language of the channel's default title and description",
	)
	updateCmd.Flags().StringVarP(
		&description, "description", "d", "", "Description of the channel",
	)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the channel")

	updateCmd.MarkFlagRequired("id")
}

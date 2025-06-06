package channel

import (
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/spf13/cobra"
)

const (
	updateShort       = "Update channel's info"
	updateLong        = "Update channel's info, such as title, description, etc"
	updateIdUsage     = "ID of the channel to update"
	updateOutputUsage = "json, yaml, or silent"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := channel.NewChannel(
			channel.WithIDs(ids),
			channel.WithCountry(country),
			channel.WithCustomUrl(customUrl),
			channel.WithDefaultLanguage(defaultLanguage),
			channel.WithDescription(description),
			channel.WithTitle(title),
			channel.WithService(nil),
		)
		c.Update(output)
	},
}

func init() {
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&country, "country", "c", "", countryUsage)
	updateCmd.Flags().StringVarP(&customUrl, "customUrl", "u", "", curlUsage)
	updateCmd.Flags().StringVarP(
		&defaultLanguage, "defaultLanguage", "l", "", dlUsage,
	)
	updateCmd.Flags().StringVarP(
		&description, "description", "d", "", descUsage,
	)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", updateOutputUsage)

	updateCmd.MarkFlagRequired("id")
}

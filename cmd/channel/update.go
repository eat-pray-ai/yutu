package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/spf13/cobra"
	"io"
)

const (
	updateShort   = "Update channel's info"
	updateLong    = "Update channel's info, such as title, description, etc"
	updateIdUsage = "ID of the channel to update"
)

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
	updateCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func update(writer io.Writer) error {
	c := channel.NewChannel(
		channel.WithIDs(ids),
		channel.WithCountry(country),
		channel.WithCustomUrl(customUrl),
		channel.WithDefaultLanguage(defaultLanguage),
		channel.WithDescription(description),
		channel.WithTitle(title),
		channel.WithService(nil),
	)

	return c.Update(output, jpath, writer)
}

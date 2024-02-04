package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "subcommand for channel manipulation",
	Long:  `subcommand for channel manipulation, which can be used to manipulate YouTube channels, such as showing channel details, updating channel details, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("channel called")
	},
}

func init() {
	rootCmd.AddCommand(channelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// channelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// channelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

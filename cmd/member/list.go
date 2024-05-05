package member

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list channel's members",
	Long:  "list channel's members' info, such as channelId, displayName, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		m := yutuber.NewMember(
			yutuber.WithMemberChannelId(memberChannelId),
		)
		m.List(parts, output)
	},
}

func init() {
	memberCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&memberChannelId, "memberChannelId", "c", "",
		"Comma separated list of channel IDs. Only data about members that are part of this list will be included",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

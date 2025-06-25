package member

import (
	"github.com/eat-pray-ai/yutu/pkg/member"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		m := member.NewMember(
			member.WithMemberChannelId(memberChannelId),
			member.WithHasAccessToLevel(hasAccessToLevel),
			member.WithMaxResults(maxResults),
			member.WithMode(mode),
			member.WithService(nil),
		)

		err := m.List(parts, output, cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func init() {
	memberCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&memberChannelId, "memberChannelId", "c", "", mcidUsage,
	)
	listCmd.Flags().StringVarP(
		&hasAccessToLevel, "hasAccessToLevel", "a", "", hatlUsage,
	)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, mrUsage)
	listCmd.Flags().StringVarP(&mode, "mode", "m", "all_current", mmUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "", outputUsage)
}

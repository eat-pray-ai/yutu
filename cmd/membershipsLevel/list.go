package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/pkg/yutuber/membershipsLevel"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list memberships levels",
	Long:  "list memberships levels' info, such as id, displayName, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		m := membershipsLevel.NewMembershipsLevel()
		m.List(parts, output)
	},
}

func init() {
	membershipsLevelCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id, snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "Output format: json or yaml",
	)
}

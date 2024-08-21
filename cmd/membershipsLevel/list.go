package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/pkg/membershipsLevel"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List memberships levels",
	Long:  "List memberships levels' info, such as id, displayName, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		m := membershipsLevel.NewMembershipsLevel(membershipsLevel.WithService(nil))
		m.List(parts, output)
	},
}

func init() {
	membershipsLevelCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id, snippet"}, "Comma separated parts",
	)
	listCmd.Flags().StringVarP(
		&output, "output", "o", "", "format: json or yaml",
	)
}

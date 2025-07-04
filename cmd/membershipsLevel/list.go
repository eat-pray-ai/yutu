package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/membershipsLevel"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	membershipsLevelCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id, snippet"}, partsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", outputUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func list(writer io.Writer) error {
	m := membershipsLevel.NewMembershipsLevel(membershipsLevel.WithService(nil))

	return m.List(parts, output, jpath, writer)
}

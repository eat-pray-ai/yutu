package channelSection

import (
	"github.com/eat-pray-ai/yutu/pkg/channelSection"
	"github.com/spf13/cobra"
	"io"
)

const (
	deleteShort    = "Delete channel sections"
	deleteLong     = "Delete channel sections by ids"
	deleteIdsUsage = "Delete the channel sections with the given ids"
)

func init() {
	channelSectionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)

	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := del(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func del(writer io.Writer) error {
	cs := channelSection.NewChannelSection(
		channelSection.WithIDs(ids),
		channelSection.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		channelSection.WithService(nil),
	)

	return cs.Delete(writer)
}

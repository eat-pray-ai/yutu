package caption

import (
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/spf13/cobra"
)

const (
	deleteShort = "Delete caption"
	deleteLong  = "Delete caption of a video"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: deleteShort,
	Long:  deleteLong,
	Run: func(cmd *cobra.Command, args []string) {
		c := caption.NewCation(
			caption.WithID(id),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithService(nil),
		)
		c.Delete()
	},
}

func init() {
	captionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", idUsage)
	deleteCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
}

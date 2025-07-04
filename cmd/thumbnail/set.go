package thumbnail

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/spf13/cobra"
	"io"
)

func init() {
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	setCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", cmd.JpUsage)

	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func set(writer io.Writer) error {
	t := thumbnail.NewThumbnail(
		thumbnail.WithFile(file),
		thumbnail.WithVideoId(videoId),
		thumbnail.WithService(nil),
	)

	return t.Set(output, jpath, writer)
}

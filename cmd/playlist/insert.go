package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/spf13/cobra"
	"io"
)

const (
	insertShort    = "Create a new playlist"
	insertLong     = "Create a new playlist, with the specified title, description, tags, etc"
	insertCidUsage = "Channel id of the playlist"
)

func init() {
	playlistCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringSliceVarP(&tags, "tags", "a", []string{}, tagsUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", insertCidUsage)
	insertCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", cmd.SilentUsage)
	insertCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", cmd.JpUsage)

	_ = insertCmd.MarkFlagRequired("title")
	_ = insertCmd.MarkFlagRequired("channel")
	_ = insertCmd.MarkFlagRequired("privacy")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := insert(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insert(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithTitle(title),
		playlist.WithDescription(description),
		playlist.WithTags(tags),
		playlist.WithLanguage(language),
		playlist.WithChannelId(channelId),
		playlist.WithPrivacy(privacy),
		playlist.WithService(nil),
	)

	return p.Insert(output, jpath, writer)
}

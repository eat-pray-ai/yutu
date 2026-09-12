// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	short         = "Manage YouTube playlists"
	long          = "Manage YouTube playlists. Use this tool to list, create, update, or delete playlists."
	titleUsage    = "Title of the playlist"
	descUsage     = "Description of the playlist"
	hlUsage       = "Return content in specified language"
	mineUsage     = "Return the playlists owned by the authenticated user"
	tagsUsage     = "Comma separated tags"
	languageUsage = "Language of the playlist"
	privacyUsage  = "public|private|unlisted"
)

var (
	ids         []string
	title       string
	description = new("")
	hl          string
	maxResults  int64
	mine        = new(false)
	tags        []string
	language    = new("")
	channelId   string
	privacy     string
	parts       []string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		utils.ResetFlags(map[string]**bool{"mine": &mine}, cmd.Flags())
		utils.ResetFlags(map[string]*[]string{"tags": &tags}, cmd.Flags())

		stringFlags := map[string]**string{
			"description": &description,
			"language":    &language,
		}
		utils.ResetFlags(stringFlags, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}

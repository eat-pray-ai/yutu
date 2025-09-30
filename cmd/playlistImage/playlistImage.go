// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short       = "Manipulate YouTube playlist images"
	long        = "List, insert, update, or delete YouTube playlist images"
	idsUsage    = "IDs of the playlist images to delete"
	heightUsage = "The image height"
	pidUsage    = "ID of the playlist this image is associated with"
	typeUsage   = "The image type (e.g., 'hero')"
	widthUsage  = "The image width"
	fileUsage   = "Path to the image file"
	parentUsage = "Return PlaylistImages for this playlist id"
)

var (
	ids        []string
	height     int64
	playlistId string
	type_      string
	width      int64
	file       string
	parent     string
	maxResults int64
	parts      []string
	output     string
	jpath      string

	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var playlistImageCmd = &cobra.Command{
	Use:   "playlistImage",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(playlistImageCmd)
}

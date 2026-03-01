// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short     = "Manage YouTube video thumbnails"
	long      = "Manage YouTube video thumbnails. Use this tool to set custom thumbnails for videos."
	fileUsage = "Path to the thumbnail file"
	vidUsage  = "ID of the video"
)

var (
	file     string
	videoId  string
	output   string
	jsonpath string
)

var thumbnailCmd = &cobra.Command{
	Use:   "thumbnail",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(thumbnailCmd)
}

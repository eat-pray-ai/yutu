// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short    = "Manage YouTube video categories"
	long     = "Manage YouTube video categories. Use this tool to list available video categories."
	idsUsage = "IDs of the video categories"
	hlUsage  = "Host language"
	rcUsage  = "Region code"
	vcURI    = "video://category/{hl}"
	vcName   = "video categories"
)

var (
	ids          []string
	hl           string
	regionCode   string
	parts        []string
	output       string
	defaultParts = []string{"id", "snippet"}
)

var videoCategoryCmd = &cobra.Command{
	Use:   "videoCategory",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoCategoryCmd)
}

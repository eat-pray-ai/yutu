// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short    = "List YouTube video categories"
	long     = "List YouTube video categories' info, such as id, title, assignable, etc."
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
	jsonpath     string
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

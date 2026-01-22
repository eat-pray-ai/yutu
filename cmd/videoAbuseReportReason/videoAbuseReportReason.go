// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short   = "List YouTube video abuse report reasons"
	long    = "List YouTube video abuse report reasons"
	hlUsage = "Host language"
)

var (
	hl       string
	parts    []string
	output   string
	jsonpath string
)

var videoAbuseReportReasonCmd = &cobra.Command{
	Use:   "videoAbuseReportReason",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(videoAbuseReportReasonCmd)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short = "Manage YouTube memberships levels"
	long  = "Manage YouTube memberships levels. Use this tool to list information about channel membership levels."
)

var (
	parts  []string
	output string
)

var membershipsLevelCmd = &cobra.Command{
	Use:   "membershipsLevel",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(membershipsLevelCmd)
}

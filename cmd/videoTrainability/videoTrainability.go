// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoTrainability

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

var videoTrainabilityCmd = &cobra.Command{
	Use:   "videoTrainability",
	Short: "Manage YouTube video trainability",
	Long:  "Manage YouTube video trainability",
}

func init() {
	cmd.RootCmd.AddCommand(videoTrainabilityCmd)
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"github.com/spf13/cobra"
)

const (
	consoleShort = "Start agent in console mode"
	consoleLong  = "Start agent in console mode"
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: consoleShort,
	Long:  consoleLong,
	Run: func(cmd *cobra.Command, args []string) {
		launcherArgs := []string{"console"}
		if streamingMode != "" {
			launcherArgs = append(launcherArgs, "-streaming_mode", streamingMode)
		}
		launch(cmd.Context(), launcherArgs)
	},
}

func init() {
	agentCmd.AddCommand(consoleCmd)
	consoleCmd.Flags().StringVarP(&streamingMode, "streamingMode", "m", "sse", "defines streaming mode: none|sse")
}

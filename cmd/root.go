// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	short = "The AI-powered toolkit that grows your YouTube channel on autopilot"
	long  = `yutu is a CLI, MCP server, and AI agent for YouTube that can automate almost all YouTube workflows.

Environment variables:
  YUTU_CREDENTIAL    Path/base64/JSON of OAuth client secret (default: client_secret.json)
  YUTU_CACHE_TOKEN   Path/base64/JSON of cached OAuth token (default: youtube.token.json)
  YUTU_ROOT          Root directory for file resolution (default: current working directory)
  YUTU_LOG_LEVEL     Log level: DEBUG, INFO, WARN, ERROR (default: INFO)`
)

var RootCmd = &cobra.Command{
	Use:   "yutu",
	Short: short,
	Long:  long,

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

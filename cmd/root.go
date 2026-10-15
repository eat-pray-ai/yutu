// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	short = "CLI, MCP server, and AI agent to automate YouTube workflows"
	long  = `yutu is a CLI, MCP server, and AI agent that automates your YouTube workflows from the terminal or any AI client.

Environment variables:
  YUTU_CREDENTIAL    Path/Base64/JSON of OAuth client secret (default: client_secret.json)
  YUTU_CACHE_TOKEN   Path/Base64/JSON of cached OAuth token (default: youtube.token.json)
  YUTU_ROOT          Root directory for file resolution (default: current working directory)
  YUTU_LOG_LEVEL     Log level: DEBUG, INFO, WARN, ERROR (default: INFO)`
)

var RootCmd = &cobra.Command{
	Use:          "yutu",
	Short:        short,
	Long:         long,
	SilenceUsage: true,

	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	cmd, err := RootCmd.ExecuteC()
	if err != nil {
		if cmd == nil {
			cmd = RootCmd
		}
		_, _ = fmt.Fprintf(
			cmd.ErrOrStderr(), "Run '%s -h' to see usage.\n", cmd.CommandPath(),
		)
		os.Exit(1)
	}
}

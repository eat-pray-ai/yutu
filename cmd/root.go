// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	short = "A fully functional MCP server and CLI for YouTube"
	long  = `yutu is a fully functional MCP server and CLI for YouTube, which can manipulate almost all YouTube resources.

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

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}

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
	long  = "yutu is a fully functional MCP server and CLI for YouTube, which can manipulate almost all YouTube resources"
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

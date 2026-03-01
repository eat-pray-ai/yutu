// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	cobramcp "github.com/eat-pray-ai/cobra-mcp"
)

var Server, mcpCmd = cobramcp.ServerAndCommand(
	&cobramcp.Config{
		Name:         "yutu",
		Version:      Version,
		Instructions: "Automate YouTube operations",
	},
)

func init() {
	RootCmd.AddCommand(mcpCmd)
}

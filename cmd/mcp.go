// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import cobramcp "github.com/eat-pray-ai/cobra-mcp"

const example = `# Start MCP server in stdio mode
yutu mcp
# Start MCP server in Streaming HTTP mode
yutu mcp --mode http --port 8216
`

var Server, mcpCmd = cobramcp.ServerAndCommand(
	&cobramcp.Config{
		Name:         "yutu",
		Version:      Version,
		Instructions: "Automate YouTube operations",
	},
)

func init() {
	mcpCmd.Example = example
	RootCmd.AddCommand(mcpCmd)
}

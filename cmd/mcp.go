// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/spf13/cobra"
)

const example = `# Start MCP server in stdio mode
yutu mcp
# Start MCP server in Streaming HTTP mode
yutu mcp --mode http --port 8216
# Start MCP server in HTTP mode with OAuth
yutu mcp --mode http --auth --port 8216
`

var mcpConfig = &cobramcp.Config{
	Name:         "yutu",
	Version:      Version,
	Instructions: "Automate YouTube operations",
}

var Server, mcpCmd = cobramcp.ServerAndCommand(mcpConfig)

var mcpAuth bool

func init() {
	mcpCmd.Example = example
	RootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().BoolVarP(
		&mcpAuth, "auth", "a", false,
		"Enable MCP OAuth authorization (HTTP mode only)",
	)

	mcpCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if mcpAuth {
			mcpConfig.Auth = &cobramcp.AuthConfig{
				TokenVerifier:        auth.GoogleTokenVerifier,
				Scopes:               auth.Scopes,
				AuthorizationServers: []string{"https://accounts.google.com"},
			}
		} else {
			port, _ := cmd.Flags().GetInt("port")
			redirectURL := fmt.Sprintf("http://localhost:%d", port)
			cmd.SetContext(common.CtxWithRedirectURL(cmd.Context(), redirectURL))
		}

		return nil
	}
}

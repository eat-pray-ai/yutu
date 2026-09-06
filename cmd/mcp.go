// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const example = `# Start MCP server in stdio mode (default)
yutu mcp

# Start MCP server in HTTP mode at http://localhost:8216/mcp
# Requires OAuth (Google account) for all MCP tool calls
yutu mcp --mode http`

var mcpConfig = &cobramcp.Config{
	Name:         "yutu",
	Version:      Version,
	ListCacheTTL: 7_200_000,
	ServerOptions: &mcp.ServerOptions{
		Instructions: "Automate YouTube operations",
	},
	HTTPOptions: &mcp.StreamableHTTPOptions{
		Stateless:                    true,
		JSONResponse:                 true,
		MaxRequestBodyBytes:          2 << 20,
		PropagateRequestCancellation: true,
	},
}

var Server, mcpCmd = cobramcp.ServerAndCommand(mcpConfig)

func init() {
	mcpCmd.Example = example
	RootCmd.AddCommand(mcpCmd)

	_ = mcpCmd.Flags().Set("host", "localhost")
	_ = mcpCmd.Flags().Set("port", "8216")
	_ = mcpCmd.Flags().MarkHidden("host")
	_ = mcpCmd.Flags().MarkHidden("port")
	_ = mcpCmd.Flags().MarkHidden("baseUrl")
	_ = mcpCmd.Flags().MarkHidden("stateless")

	mcpCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		mode, _ := cmd.Flags().GetString("mode")
		if mode == "http" {
			mcpConfig.Auth = &cobramcp.AuthConfig{
				TokenVerifier:        auth.GoogleTokenVerifier,
				Scopes:               auth.Scopes,
				AuthorizationServers: []string{"https://accounts.google.com"},
			}
		}

		return nil
	}
}

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

# Start MCP server in HTTP mode (OAuth required, binds 127.0.0.1)
yutu mcp --mode http --port 8216

# Bind to all interfaces (e.g. container deployment)
yutu mcp --mode http --host 0.0.0.0 --port 8216

# Behind a reverse proxy with a public base URL
yutu mcp --mode http --baseUrl https://mcp.example.com`

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

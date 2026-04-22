// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
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

var (
	mcpAuth bool
	baseURL string
)

func init() {
	mcpCmd.Example = example
	RootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().BoolVarP(
		&mcpAuth, "auth", "a", false,
		"Enable MCP OAuth authorization (HTTP mode only)",
	)
	mcpCmd.Flags().StringVarP(
		&baseURL, "baseUrl", "b", "",
		"Base URL for the MCP server (default http://localhost:<port>)",
	)

	mcpCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")

		if mcpAuth {
			if baseURL == "" {
				baseURL = fmt.Sprintf("http://localhost:%d", port)
			}
			mcpConfig.Auth = &cobramcp.AuthConfig{
				ResourceMetadata: &oauthex.ProtectedResourceMetadata{
					Resource:             baseURL + "/mcp",
					AuthorizationServers: []string{"https://accounts.google.com"},
					ScopesSupported:      auth.Scopes,
				},
				ResourceMetadataURL: baseURL + "/.well-known/oauth-protected-resource",
				TokenVerifier:       auth.GoogleTokenVerifier,
				Scopes:              auth.Scopes,
			}
		} else {
			common.RedirectURL = fmt.Sprintf("http://localhost:%d", port)
		}

		return nil
	}
}

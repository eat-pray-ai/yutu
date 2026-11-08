// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"iter"
	"os/exec"
	"slices"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	example = `# Start MCP server in stdio mode (default)
yutu mcp

# Start MCP server in HTTP mode at http://localhost:8216/mcp
# Requires OAuth (Google account) for all MCP tool calls
yutu mcp --mode http`
	setupShort   = "Setup yutu MCP in AI Agents"
	setupLong    = "Setup yutu as a global stdio MCP server for the selected AI Agents."
	setupExample = `# Setup MCP for multiple AI Agents
yutu mcp setup -a claude,opencode -c client_secret.json -t youtube.token.json

# Setup MCP for Antigravity
yutu mcp setup -a agy -c client_secret.json -t youtube.token.json`
	agentsUsage     = "Comma-separated AI Agents: claude,opencode,pi,codex,agy"
	credentialUsage = "Path/Base64/JSON of OAuth client secret"
	cacheTokenUsage = "Path/Base64/JSON of cached OAuth token"
)

var (
	agentNames = []string{"claude", "opencode", "pi", "codex", "agy"}
	agents     []string
)

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

func SetToolExecutionHook(hook cobramcp.ToolExecutionHook) {
	cobramcp.SetDefaultToolExecutionHook(hook)
}

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

	mcpCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringSliceVarP(&agents, "agents", "a", nil, agentsUsage)
	setupCmd.Flags().StringVarP(&credential, "credential", "c", "client_secret.json", credentialUsage)
	setupCmd.Flags().StringVarP(&cacheToken, "cacheToken", "t", "youtube.token.json", cacheTokenUsage)
	_ = setupCmd.MarkFlagRequired("agents")
}

var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   setupShort,
	Long:    setupLong,
	Example: setupExample,
	PreRunE: func(c *cobra.Command, _ []string) (errs error) {
		var err error
		credential, err = utils.ResolveCredential(c, "credential", "YUTU_CREDENTIAL", credential)
		errs = errors.Join(errs, err)
		cacheToken, err = utils.ResolveCredential(c, "cacheToken", "YUTU_CACHE_TOKEN", cacheToken)
		errs = errors.Join(errs, err)
		return errs
	},
	RunE: func(c *cobra.Command, _ []string) (errs error) {
		c.SilenceErrors = true
		for agent, err := range parseAgents(agents) {
			if err != nil {
				_, _ = fmt.Fprintf(c.ErrOrStderr(), "Error: %v\n", err)
				errs = errors.Join(errs, err)
				continue
			}
			_, _ = fmt.Fprintf(c.OutOrStdout(), "Setup %s... ", agent)
			if err := runAgentCommand(c.Context(), agent, buildAgentArgs(agent)); err != nil {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "Failed: %v\n", err)
				errs = errors.Join(errs, fmt.Errorf("%s: %w", agent, err))
				continue
			}
			_, _ = fmt.Fprintln(c.OutOrStdout(), "Done")
		}
		return errs
	},
}

func buildAgentArgs(agent string) []string {
	var env []string
	if credential != "" {
		env = append(env, "--env", "YUTU_CREDENTIAL="+credential)
	}
	if cacheToken != "" {
		env = append(env, "--env", "YUTU_CACHE_TOKEN="+cacheToken)
	}

	args := []string{"mcp", "add"}
	switch agent {
	case "claude":
		args = append(args, "yutu", "--scope", "user", "--transport", "stdio")
		args = append(args, env...)
	case "opencode":
		args = append(args, "yutu", "--global")
		args = append(args, env...)
	case "pi", "codex":
		args = append(args, "yutu")
		args = append(args, env...)
	case "agy":
		args = append(args, "--type", "stdio")
		args = append(args, env...)
		args = append(args, "yutu")
	}
	return append(args, "--", "yutu", "mcp", "-m", "stdio")
}

func parseAgents(values []string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		var seen []string
		for _, value := range values {
			name := strings.ToLower(strings.TrimSpace(value))
			if !slices.Contains(agentNames, name) {
				err := fmt.Errorf("unsupported Agent %q: choose %s", name, strings.Join(agentNames, ", "))
				if !yield("", err) {
					return
				}
				continue
			}
			if slices.Contains(seen, name) {
				continue
			}
			seen = append(seen, name)
			if !yield(name, nil) {
				return
			}
		}
		if len(values) == 0 {
			yield("", errors.New("--agents requires at least one AI Agent"))
		}
	}
}

func runAgentCommand(ctx context.Context, binary string, args []string) error {
	if _, err := exec.LookPath(binary); err != nil {
		return errors.New("not found on PATH; install it and rerun setup")
	}
	out, err := exec.CommandContext(ctx, binary, args...).CombinedOutput()
	if err != nil {
		return errors.New(cmp.Or(strings.TrimSpace(string(out)), "check the Agent's MCP settings"))
	}
	return nil
}

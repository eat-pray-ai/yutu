// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

const (
	agentShort = "Start agent to automate YouTube workflows"
	agentLong  = "Start agent to automate YouTube workflows"
)

var launcherArgs string

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: agentShort,
	Long:  agentLong,
	Run: func(cmd *cobra.Command, args []string) {
		if launcherArgs == "" {
			_ = cmd.Help()
			return
		}
		launch(cmd.Context(), cmd.OutOrStdout(), strings.Fields(launcherArgs))
	},
}

func init() {
	cmd.RootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVarP(
		&launcherArgs, "args", "a", "",
		"launcher arguments as a single string, e.g. 'console -streaming_mode sse'",
	)
}

func launch(ctx context.Context, writer io.Writer, args []string) {
	m, err := gemini.NewModel(
		ctx, os.Getenv("YUTU_AGENT_MODEL"), &genai.ClientConfig{
			APIKey:  os.Getenv("YUTU_LLM_API_KEY"),
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create Gemini model", "error", err)
		os.Exit(1)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	_, err = cmd.Server.Connect(ctx, serverTransport, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to MCP server", "error", err)
		os.Exit(1)
	}

	mcpToolSet, err := mcptoolset.New(
		mcptoolset.Config{
			Transport: clientTransport,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create MCP tool set", "error", err)
	}

	orchestrator, err := buildOrchestrator(m, mcpToolSet)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build orchestrator agent", "error", err)
		os.Exit(1)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(orchestrator),
	}
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, args); err != nil {
		slog.ErrorContext(ctx, "failed to launch agent", "error", err)
		_, _ = fmt.Fprintln(writer, l.CommandLineSyntax())
		os.Exit(1)
	}
}

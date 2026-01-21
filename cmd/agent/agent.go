// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"time"

	rootCmd "github.com/eat-pray-ai/yutu/cmd"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

const (
	agentShort = "Start agent to automate YouTube workflows"
	agentLong  = "Start agent to automate YouTube workflows"
)

var (
	//go:embed INSTRUCTION.md
	instruction string

	streamingMode    string
	idleTimeout      time.Duration
	port             int
	readTimeout      time.Duration
	writeTimeout     time.Duration
	sublaunchers     []string
	sseWriteTimeout  time.Duration
	webuiAddress     string
	a2aAgentURL      string
	apiServerAddress string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: agentShort,
	Long:  agentLong,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		envInstruction, ok := os.LookupEnv("YUTU_AGENT_INSTRUCTION")
		if ok && envInstruction != "" {
			instruction = envInstruction
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.RootCmd.AddCommand(agentCmd)
}

func launch(ctx context.Context, args []string) {
	m, err := gemini.NewModel(
		ctx, os.Getenv("GEMINI_MODEL"), &genai.ClientConfig{
			APIKey:  os.Getenv("GEMINI_API_KEY"),
			Backend: genai.BackendGeminiAPI,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create Gemini model", "error", err)
		os.Exit(1)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	_, err = rootCmd.Server.Connect(ctx, serverTransport, nil)
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

	a, err := llmagent.New(
		llmagent.Config{
			Model:       m,
			Name:        "YouTube Copilot",
			Instruction: instruction,
			Toolsets:    []tool.Toolset{mcpToolSet},
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create LLM agent", "error", err)
		os.Exit(1)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, args); err != nil {
		slog.ErrorContext(
			ctx, "failed to launch agent",
			"launch", l.CommandLineSyntax(), "error", err,
		)
		os.Exit(1)
	}
}

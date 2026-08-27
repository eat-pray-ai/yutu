// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/geminitool"
	"google.golang.org/adk/v2/tool/mcptoolset"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
	"google.golang.org/genai"
)

const (
	short   = "Start an agent to automate YouTube workflows"
	long    = "Start an agent to automate YouTube workflows."
	example = `# console mode
yutu agent --args "console" --model "google:gemini-3.7-flash" --api-key "YOUR_KEY"
# web mode
yutu agent --args "web api a2a webui" --model "google:gemini-3.7-flash" --api-key "YOUR_KEY"`
	argsUsage        = "Launcher arguments as a single string"
	modelUsage       = "Model in provider:modelName format"
	apiKeyUsage      = "API key for the model provider"
	instructionUsage = "Override the built-in agent instruction"

	agentDescription       = "YouTube growth strategist and workflow assistant — retrieves, creates, updates, and deletes YouTube content."
	errInvalidModelSpec    = "invalid model spec %q: expected provider:modelName, e.g. google:gemini-3.7-flash"
	errUnsupportedProvider = "unsupported provider %q: only \"google\" is supported"
)

var (
	launcherArgs        string
	modelSpec           string
	apiKey              string
	instructionOverride string

	//go:embed INSTRUCTION.md
	instruction string

	//go:embed skills
	skillsFS embed.FS
)

var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   short,
	Long:    long,
	Example: example,
	Run: func(cmd *cobra.Command, _ []string) {
		if launcherArgs == "" {
			_, _ = fmt.Fprintln(
				cmd.OutOrStdout(), full.NewLauncher().CommandLineSyntax(),
			)
			return
		}
		launch(cmd.Context(), cmd.OutOrStdout(), strings.Fields(launcherArgs))
	},
}

func init() {
	cmd.RootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVarP(&launcherArgs, "args", "a", "console", argsUsage)
	agentCmd.Flags().StringVarP(
		&modelSpec, "model", "m", "google:gemini-3.7-flash", modelUsage,
	)
	agentCmd.Flags().StringVar(&apiKey, "api-key", "", apiKeyUsage)
	agentCmd.Flags().StringVarP(
		&instructionOverride, "instruction", "i", "", instructionUsage,
	)
}

func newModel(ctx context.Context, spec string) (model.LLM, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf(errInvalidModelSpec, spec)
	}

	provider, modelName := parts[0], parts[1]
	switch provider {
	case "google":
		return gemini.NewModel(
			ctx, modelName, &genai.ClientConfig{
				APIKey:  apiKey,
				Backend: genai.BackendGeminiAPI,
			},
		)
	default:
		return nil, fmt.Errorf(errUnsupportedProvider, provider)
	}
}

func newSkillToolset(ctx context.Context) (tool.Toolset, error) {
	source := skill.NewFileSystemSource(skillsFS)
	source, _, err := skill.WithCompletePreloadSource(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("preload skills: %w", err)
	}
	return skilltoolset.New(ctx, skilltoolset.Config{Source: source})
}

func buildAgent(
	m model.LLM, mcpToolSet, skillToolset tool.Toolset,
) (agent.Agent, error) {
	return llmagent.New(
		llmagent.Config{
			Name:        "Miffy",
			Model:       m,
			Description: agentDescription,
			Instruction: instruction,
			Tools:       []tool.Tool{geminitool.GoogleSearch{}},
			Toolsets:    []tool.Toolset{mcpToolSet, skillToolset},
		},
	)
}

func launch(ctx context.Context, writer io.Writer, args []string) {
	if instructionOverride != "" {
		instruction = instructionOverride
	}

	m, err := newModel(ctx, modelSpec)
	if err != nil {
		slog.ErrorContext(ctx, "model configuration error", "error", err)
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
		os.Exit(1)
	}

	skillToolset, err := newSkillToolset(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create skill toolset", "error", err)
		os.Exit(1)
	}

	miffy, err := buildAgent(m, mcpToolSet, skillToolset)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build agent", "error", err)
		os.Exit(1)
	}

	sessionService := session.InMemoryService()
	memoryService := memory.InMemoryService()

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(miffy),
		SessionService: sessionService,
		MemoryService:  memoryService,
	}
	l := full.NewLauncher()
	if err := l.Execute(ctx, config, args); err != nil {
		slog.ErrorContext(ctx, "failed to launch agent", "error", err)
		_, _ = fmt.Fprintln(writer, l.CommandLineSyntax())
		os.Exit(1)
	}
}

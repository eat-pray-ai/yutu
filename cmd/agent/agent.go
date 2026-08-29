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
	"google.golang.org/adk/v2/model/openaimodel"
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
	example = `# Gemini (default)
yutu agent --provider google --model gemini-3.7-flash --api-key "YOUR_KEY"
# OpenAI
yutu agent --provider openai --model gpt-5.6-terra
# OpenAI-compatible endpoint (DeepSeek, Ollama, vLLM, etc.)
yutu agent --provider openai-compatible --model deepseek-v4-flash --base-url "https://api.deepseek.com/"
# Custom Gemini endpoint
yutu agent --provider google --model gemini-3.7-flash --base-url "https://custom-endpoint.example.com/"`
	argsUsage        = "Launcher arguments as a single string"
	providerUsage    = "LLM provider (google, openai, openai-compatible)"
	modelUsage       = "Model name"
	apiKeyUsage      = "API key for the model provider"
	baseURLUsage     = "Base URL for the model provider's API endpoint"
	instructionUsage = "Override the built-in agent instruction"
	agentDescription = "YouTube growth strategist and workflow assistant — retrieve, create, update, and delete YouTube content."

	errUnsupportedProvider = "unsupported provider %q: supported providers are google, openai, openai-compatible"
	errModelConfig         = "model configuration error"
	errMCPConnect          = "failed to connect to MCP server"
	errMCPToolSet          = "failed to create MCP tool set"
	errSkillToolset        = "failed to create skill toolset"
	errBuildAgent          = "failed to build agent"
	errLaunchAgent         = "failed to launch agent"
)

var (
	launcherArgs string
	provider     string
	modelName    string
	apiKey       string
	baseURL      string
	inst         string

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
		&provider, "provider", "p", "google", providerUsage,
	)
	agentCmd.Flags().StringVarP(
		&modelName, "model", "m", "gemini-3.7-flash", modelUsage,
	)
	agentCmd.Flags().StringVar(&apiKey, "api-key", "", apiKeyUsage)
	agentCmd.Flags().StringVar(&baseURL, "base-url", "", baseURLUsage)
	agentCmd.Flags().StringVarP(&inst, "instruction", "i", "", instructionUsage)
}

func newModel(ctx context.Context) (model.LLM, error) {
	switch provider {
	case "google":
		cfg := &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		}
		if baseURL != "" {
			cfg.HTTPOptions.BaseURL = baseURL
		}
		return gemini.NewModel(ctx, modelName, cfg)
	case "openai", "openai-compatible":
		return openaimodel.NewModel(
			ctx, modelName, &openaimodel.ClientConfig{
				APIKey:  apiKey,
				BaseURL: baseURL,
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
	var tools []tool.Tool
	if provider == "google" {
		tools = append(tools, geminitool.GoogleSearch{})
	}
	return llmagent.New(
		llmagent.Config{
			Name:        "Miffy",
			Model:       m,
			Description: agentDescription,
			Instruction: instruction,
			Tools:       tools,
			Toolsets:    []tool.Toolset{mcpToolSet, skillToolset},
		},
	)
}

func launch(ctx context.Context, writer io.Writer, args []string) {
	if inst != "" {
		instruction = inst
	}

	m, err := newModel(ctx)
	if err != nil {
		slog.ErrorContext(ctx, errModelConfig, "error", err)
		os.Exit(1)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	_, err = cmd.Server.Connect(ctx, serverTransport, nil)
	if err != nil {
		slog.ErrorContext(ctx, errMCPConnect, "error", err)
		os.Exit(1)
	}

	mcpToolSet, err := mcptoolset.New(
		mcptoolset.Config{
			Transport: clientTransport,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, errMCPToolSet, "error", err)
		os.Exit(1)
	}

	skillToolset, err := newSkillToolset(ctx)
	if err != nil {
		slog.ErrorContext(ctx, errSkillToolset, "error", err)
		os.Exit(1)
	}

	miffy, err := buildAgent(m, mcpToolSet, skillToolset)
	if err != nil {
		slog.ErrorContext(ctx, errBuildAgent, "error", err)
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
		slog.ErrorContext(ctx, errLaunchAgent, "error", err)
		_, _ = fmt.Fprintln(writer, l.CommandLineSyntax())
		os.Exit(1)
	}
}

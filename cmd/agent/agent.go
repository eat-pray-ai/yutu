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
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

const (
	short = "Start an agent to automate YouTube workflows"
	long  = `Start an agent to automate YouTube workflows.

Environment variables:
  YUTU_ADVANCED_MODEL          Model for orchestrator agent (format: provider:modelName, e.g. google:gemini-3.1-pro-preview)
  YUTU_LITE_MODEL              Model for sub-agents (format: provider:modelName, e.g. google:gemini-3-flash-preview)
  YUTU_LLM_API_KEY             API key for the model provider
  GOOGLE_GEMINI_BASE_URL       Base URL for Gemini API (optional)
  YUTU_AGENT_INSTRUCTION       Custom instruction for orchestrator agent (optional)
  YUTU_RETRIEVAL_INSTRUCTION   Custom instruction for retrieval agent (optional)
  YUTU_MODIFIER_INSTRUCTION    Custom instruction for modifier agent (optional)
  YUTU_DESTROYER_INSTRUCTION   Custom instruction for destroyer agent (optional)

At least one of YUTU_ADVANCED_MODEL or YUTU_LITE_MODEL must be set.
If only one is set, the other defaults to the same value.`
	example = `# console mode
yutu agent --args "console"
# web mode with three sub-launchers: api, a2a and webui
yutu agent --args "web api a2a webui"`
)

var launcherArgs string

var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   short,
	Long:    long,
	Example: example,
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

func resolveModels(ctx context.Context) (
	advancedModel, liteModel model.LLM, err error,
) {
	advanced := os.Getenv("YUTU_ADVANCED_MODEL")
	lite := os.Getenv("YUTU_LITE_MODEL")

	if advanced == "" && lite == "" {
		return nil, nil, fmt.Errorf(
			"at least one of YUTU_ADVANCED_MODEL or YUTU_LITE_MODEL must be set (format: provider:modelName)",
		)
	}
	if advanced == "" {
		advanced = lite
	}
	if lite == "" {
		lite = advanced
	}

	advancedModel, err = newModel(ctx, advanced)
	if err != nil {
		return nil, nil, fmt.Errorf("YUTU_ADVANCED_MODEL: %w", err)
	}

	if lite == advanced {
		return advancedModel, advancedModel, nil
	}

	liteModel, err = newModel(ctx, lite)
	if err != nil {
		return nil, nil, fmt.Errorf("YUTU_LITE_MODEL: %w", err)
	}

	return advancedModel, liteModel, nil
}

func newModel(ctx context.Context, spec string) (model.LLM, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf(
			"invalid model spec %q: expected format provider:modelName (e.g. google:gemini-3.1-pro-preview)",
			spec,
		)
	}

	provider, modelName := parts[0], parts[1]
	switch provider {
	case "google":
		return gemini.NewModel(
			ctx, modelName, &genai.ClientConfig{
				APIKey:  os.Getenv("YUTU_LLM_API_KEY"),
				Backend: genai.BackendGeminiAPI,
			},
		)
	default:
		return nil, fmt.Errorf(
			"unsupported provider %q: only \"google\" is supported", provider,
		)
	}
}

func launch(ctx context.Context, writer io.Writer, args []string) {
	advancedModel, liteModel, err := resolveModels(ctx)
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
			// RequireConfirmationProvider: requireConfirmation,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create MCP tool set", "error", err)
	}

	orchestrator, err := buildOrchestrator(advancedModel, liteModel, mcpToolSet)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build orchestrator agent", "error", err)
		os.Exit(1)
	}

	sessionService := session.InMemoryService()
	memoryService := memory.InMemoryService()

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(orchestrator),
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

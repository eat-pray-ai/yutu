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
	pkgagent "github.com/eat-pray-ai/yutu/pkg/agent"
	"github.com/spf13/cobra"
	"google.golang.org/adk/v2/cmd/launcher/full"
)

const (
	short   = "Start an Agent to automate YouTube workflows"
	long    = "Start an Agent to automate YouTube workflows."
	example = `  # Gemini(compatible) (default)
  yutu agent --provider google --model gemini-4-hydrogen --apiKey "YOUR_KEY"
  yutu agent --provider google --model gemini-4-carbon --apiKey "YOUR_KEY" --baseUrl "https://custom.endpoint.com"
  # OpenAI(compatible) (OpenAI, DeepSeek, Ollama, vLLM, etc.)
  yutu agent --provider openai --model gpt-6-astra --apiKey "YOUR_KEY"
  yutu agent --provider openai --model deepseek-flash --apiKey "YOUR_KEY" --baseUrl "https://api.deepseek.com"
  # Anthropic(compatible)
  yutu agent --provider anthropic --model claude-fable-5-2 --apiKey "YOUR_KEY"
  yutu agent --provider anthropic --model claude-opus-5-5 --apiKey "YOUR_KEY" --baseUrl "https://custom.endpoint.com"`
	argsUsage        = "Launcher arguments as a single string"
	providerUsage    = "LLM provider: google, openai or anthropic"
	modelUsage       = "Model name"
	apiKeyUsage      = "API key for the model provider"
	baseURLUsage     = "Base URL for the model provider's API endpoint"
	instructionUsage = "Override the built-in Agent instruction"

	errInitAgent   = "failed to initialize Agent"
	errLaunchAgent = "failed to launch Agent"
)

func init() {
	cmd.RootCmd.AddCommand(newAgentCmd())
}

type agentOpts struct {
	launcherArgs string
	provider     string
	modelName    string
	apiKey       string
	baseURL      string
	inst         string
}

func (o *agentOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.launcherArgs, "args", "a", "console", argsUsage)
	c.Flags().StringVarP(
		&o.provider, "provider", "p", pkgagent.DefaultProvider, providerUsage,
	)
	c.Flags().StringVarP(
		&o.modelName, "model", "m", pkgagent.DefaultModel, modelUsage,
	)
	c.Flags().StringVarP(&o.apiKey, "apiKey", "k", "", apiKeyUsage)
	c.Flags().StringVarP(&o.baseURL, "baseUrl", "u", "", baseURLUsage)
	c.Flags().StringVarP(&o.inst, "instruction", "i", "", instructionUsage)
}

func (o *agentOpts) launch(ctx context.Context, writer io.Writer, args []string) {
	instance, err := pkgagent.New(
		ctx, pkgagent.Config{
			Provider:    o.provider,
			Model:       o.modelName,
			APIKey:      o.apiKey,
			BaseURL:     o.baseURL,
			Instruction: o.inst,
			MCPServer:   cmd.Server,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, errInitAgent, "error", err)
		os.Exit(1)
	}
	defer func() { _ = instance.Close() }()

	l := full.NewLauncher()
	if err := l.Execute(ctx, instance.LauncherConfig(), args); err != nil {
		slog.ErrorContext(ctx, errLaunchAgent, "error", err)
		_, _ = fmt.Fprintln(writer, l.CommandLineSyntax())
		os.Exit(1)
	}
}

func newAgentCmd() *cobra.Command {
	var opts agentOpts

	c := &cobra.Command{
		Use:     "agent",
		Short:   short,
		Long:    long,
		Example: example,
		Run: func(c *cobra.Command, _ []string) {
			if opts.launcherArgs == "" {
				_, _ = fmt.Fprintln(
					c.OutOrStdout(), full.NewLauncher().CommandLineSyntax(),
				)
				return
			}
			opts.launch(c.Context(), c.OutOrStdout(), strings.Fields(opts.launcherArgs))
		},
	}

	opts.bindFlags(c)
	return c
}

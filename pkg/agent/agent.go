// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"iter"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/model/anthropicmodel"
	"google.golang.org/adk/v2/model/gemini"
	"google.golang.org/adk/v2/model/openaimodel"
	"google.golang.org/adk/v2/runner"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/geminitool"
	"google.golang.org/adk/v2/tool/mcptoolset"
	"google.golang.org/adk/v2/tool/skilltoolset"
	"google.golang.org/adk/v2/tool/skilltoolset/skill"
	"google.golang.org/genai"
)

const (
	Name            = "Miffy"
	DefaultProvider = "google"
	DefaultModel    = "gemini-4-hydrogen"
	Description     = "Miffy, your YouTube autopilot — orchestrates uploads, SEO, comments, playlists, and live studio operations."
)

var (
	//go:embed INSTRUCTION.md
	defaultInstruction string

	//go:embed skills
	skillsFS embed.FS

	ErrMissingMCPServer = errors.New("MCP server is required")
)

type Config struct {
	Provider    string
	Model       string
	APIKey      string
	BaseURL     string
	Instruction string
	MCPServer   *mcp.Server
}

type Instance struct {
	Agent          adkagent.Agent
	Runner         *runner.Runner
	SessionService session.Service
	MemoryService  memory.Service
	serverSession  *mcp.ServerSession
}

func New(ctx context.Context, cfg Config) (*Instance, error) {
	if cfg.MCPServer == nil {
		return nil, ErrMissingMCPServer
	}
	if cfg.Provider == "" {
		cfg.Provider = DefaultProvider
	}
	if cfg.Model == "" {
		cfg.Model = DefaultModel
	}
	if cfg.Instruction == "" {
		cfg.Instruction = defaultInstruction
	}

	m, err := newModel(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("model configuration error: %w", err)
	}

	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := cfg.MCPServer.Connect(ctx, serverTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP server: %w", err)
	}

	mcpToolSet, err := mcptoolset.New(
		mcptoolset.Config{
			Transport: clientTransport,
		},
	)
	if err != nil {
		_ = serverSession.Close()
		return nil, fmt.Errorf("failed to create MCP toolset: %w", err)
	}

	skillToolset, err := newSkillToolset(ctx)
	if err != nil {
		_ = serverSession.Close()
		return nil, fmt.Errorf("failed to create skill toolset: %w", err)
	}

	miffy, err := buildAgent(cfg, m, mcpToolSet, skillToolset)
	if err != nil {
		_ = serverSession.Close()
		return nil, fmt.Errorf("failed to build Agent: %w", err)
	}

	sessionService := session.InMemoryService()
	memoryService := memory.InMemoryService()

	r, err := runner.New(
		runner.Config{
			AppName:           Name,
			Agent:             miffy,
			SessionService:    sessionService,
			MemoryService:     memoryService,
			AutoCreateSession: true,
		},
	)
	if err != nil {
		_ = serverSession.Close()
		return nil, fmt.Errorf("failed to create Agent runner: %w", err)
	}

	return &Instance{
		Agent:          miffy,
		Runner:         r,
		SessionService: sessionService,
		MemoryService:  memoryService,
		serverSession:  serverSession,
	}, nil
}

func (i *Instance) Run(
	ctx context.Context,
	userID, sessionID string,
	msg *genai.Content,
	cfg adkagent.RunConfig,
	opts ...runner.RunOption,
) iter.Seq2[*session.Event, error] {
	return i.Runner.Run(ctx, userID, sessionID, msg, cfg, opts...)
}

func (i *Instance) Send(
	ctx context.Context,
	userID, sessionID, text string,
	opts ...runner.RunOption,
) iter.Seq2[*session.Event, error] {
	msg := genai.NewContentFromText(text, genai.RoleUser)
	return i.Runner.Run(
		ctx, userID, sessionID, msg,
		adkagent.RunConfig{StreamingMode: adkagent.StreamingModeSSE},
		opts...,
	)
}

func (i *Instance) LauncherConfig() *launcher.Config {
	return &launcher.Config{
		AgentLoader:    adkagent.NewSingleLoader(i.Agent),
		SessionService: i.SessionService,
		MemoryService:  i.MemoryService,
	}
}

func (i *Instance) Close() error {
	if i == nil || i.serverSession == nil {
		return nil
	}
	return i.serverSession.Close()
}

func newModel(ctx context.Context, cfg Config) (model.LLM, error) {
	switch cfg.Provider {
	case "google":
		clientCfg := &genai.ClientConfig{
			APIKey:  cfg.APIKey,
			Backend: genai.BackendGeminiAPI,
		}
		if cfg.BaseURL != "" {
			clientCfg.HTTPOptions.BaseURL = cfg.BaseURL
		}
		return gemini.NewModel(ctx, cfg.Model, clientCfg)
	case "openai":
		return openaimodel.NewModel(
			ctx, cfg.Model, &openaimodel.ClientConfig{
				APIKey:  cfg.APIKey,
				BaseURL: cfg.BaseURL,
			},
		)
	case "anthropic":
		return anthropicmodel.NewModel(
			ctx, cfg.Model, anthropicmodel.ClientConfig{
				APIKey:  cfg.APIKey,
				BaseURL: cfg.BaseURL,
			},
		)
	default:
		return nil, fmt.Errorf(
			"unsupported provider %q: supported providers are google, openai and anthropic",
			cfg.Provider,
		)
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
	cfg Config, m model.LLM, mcpToolSet, skillToolset tool.Toolset,
) (adkagent.Agent, error) {
	var tools []tool.Tool
	if cfg.Provider == "google" && cfg.BaseURL == "" {
		tools = append(tools, geminitool.GoogleSearch{})
	}
	return llmagent.New(
		llmagent.Config{
			Name:        Name,
			Model:       m,
			Description: Description,
			Instruction: cfg.Instruction,
			Tools:       tools,
			Toolsets:    []tool.Toolset{mcpToolSet, skillToolset},
		},
	)
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"errors"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newTestMCPServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
}

func TestNew_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr error
		wantSub string
	}{
		{
			name:    "missing MCPServer",
			cfg:     Config{Provider: "google", APIKey: "test-key"},
			wantErr: ErrMissingMCPServer,
		},
		{
			name: "unsupported provider",
			cfg: Config{
				Provider:  "invalid",
				APIKey:    "test-key",
				MCPServer: newTestMCPServer(),
			},
			wantSub: "unsupported provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst, err := New(t.Context(), tt.cfg)
			if err == nil {
				_ = inst.Close()
				t.Fatal("New() error = nil, want error")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantSub != "" && !strings.Contains(err.Error(), tt.wantSub) {
				t.Errorf("New() error = %q, want substring %q", err.Error(), tt.wantSub)
			}
		})
	}
}

func TestNew_Providers(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default google provider",
			cfg: Config{
				APIKey:    "test-key",
				MCPServer: newTestMCPServer(),
			},
		},
		{
			name: "google with custom base url and instruction",
			cfg: Config{
				Provider:    "google",
				Model:       "gemini-4-carbon",
				APIKey:      "test-key",
				BaseURL:     "https://example.com",
				Instruction: "custom instruction",
				MCPServer:   newTestMCPServer(),
			},
		},
		{
			name: "openai provider",
			cfg: Config{
				Provider:  "openai",
				Model:     "gpt-6-astra",
				APIKey:    "test-key",
				BaseURL:   "https://example.com/v1",
				MCPServer: newTestMCPServer(),
			},
		},
		{
			name: "anthropic provider",
			cfg: Config{
				Provider:  "anthropic",
				Model:     "claude-fable-5-2",
				APIKey:    "test-key",
				BaseURL:   "https://example.com",
				MCPServer: newTestMCPServer(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst, err := New(t.Context(), tt.cfg)
			if err != nil {
				t.Fatalf("New() unexpected error = %v", err)
			}
			defer func() {
				if err := inst.Close(); err != nil {
					t.Errorf("Close() error = %v", err)
				}
			}()

			if inst.Agent == nil {
				t.Fatal("Agent is nil")
			}
			if inst.Agent.Name() != Name {
				t.Errorf("Agent.Name() = %q, want %q", inst.Agent.Name(), Name)
			}
			if inst.Runner == nil {
				t.Fatal("Runner is nil")
			}
			if inst.SessionService == nil {
				t.Fatal("SessionService is nil")
			}
			if inst.MemoryService == nil {
				t.Fatal("MemoryService is nil")
			}
			lcfg := inst.LauncherConfig()
			if lcfg == nil || lcfg.AgentLoader == nil {
				t.Fatal("LauncherConfig() returned incomplete config")
			}
		})
	}
}

func TestInstance_CloseNil(t *testing.T) {
	var inst *Instance
	if err := inst.Close(); err != nil {
		t.Errorf("(*Instance)(nil).Close() = %v, want nil", err)
	}
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	_ "embed"
	"fmt"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/geminitool"
)

var (
	//go:embed INSTRUCTION_ORCHESTRATOR.md
	orchestratorInstruction string

	//go:embed INSTRUCTION_RETRIEVAL.md
	retrievalInstruction string

	//go:embed INSTRUCTION_MODIFIER.md
	modifierInstruction string

	//go:embed INSTRUCTION_DESTROYER.md
	destroyerInstruction string
)

type agentDef struct {
	name        string
	description string
	instruction *string
	envKey      string
	toolNames   []string
}

var agentDefs = map[string]agentDef{
	"Nina": {
		name:        "YouTube Copilot",
		description: "Orchestrates YouTube workflows by delegating to specialized agents.",
		instruction: &orchestratorInstruction,
		envKey:      "YUTU_AGENT_INSTRUCTION",
	},
	"Aagje": {
		name:        "Retrieval",
		description: "Retrieves information from YouTube and the internet, such as listing videos, channels, playlists, comments, and more.",
		instruction: &retrievalInstruction,
		envKey:      "YUTU_RETRIEVAL_INSTRUCTION",
		toolNames: []string{
			"activity-list",
			"caption-download",
			"caption-list",
			"channel-list",
			"channelSection-list",
			"comment-list",
			"commentThread-list",
			"member-list",
			"membershipsLevel-list",
			"playlist-list",
			"playlistImage-list",
			"playlistItem-list",
			"search-list",
			"subscription-list",
			"superChatEvent-list",
			"video-getRating",
			"video-list",
			"videoAbuseReportReason-list",
		},
	},
	"Knorretje": {
		name:        "Modifier",
		description: "Creates and updates YouTube content, such as uploading videos, creating playlists, updating metadata, posting comments, and more.",
		instruction: &modifierInstruction,
		envKey:      "YUTU_MODIFIER_INSTRUCTION",
		toolNames: []string{
			"caption-insert",
			"caption-update",
			"channel-update",
			"channelBanner-insert",
			"comment-insert",
			"comment-markAsSpam",
			"comment-setModerationStatus",
			"comment-update",
			"commentThread-insert",
			"playlist-insert",
			"playlist-update",
			"playlistImage-insert",
			"playlistImage-update",
			"playlistItem-insert",
			"playlistItem-update",
			"subscription-insert",
			"thumbnail-set",
			"video-insert",
			"video-rate",
			"video-reportAbuse",
			"video-update",
			"watermark-set",
		},
	},
	"Daan": {
		name:        "Destroyer",
		description: "Deletes YouTube content, such as deleting videos, playlists, comments, captions, or subscriptions. Handles destructive operations that require extra caution.",
		instruction: &destroyerInstruction,
		envKey:      "YUTU_DESTROYER_INSTRUCTION",
		toolNames: []string{
			"caption-delete",
			"channelSection-delete",
			"comment-delete",
			"playlist-delete",
			"playlistImage-delete",
			"playlistItem-delete",
			"subscription-delete",
			"video-delete",
			"watermark-unset",
		},
	},
}

func init() {
	for _, def := range agentDefs {
		if v, ok := os.LookupEnv(def.envKey); ok && v != "" {
			*def.instruction = v
		}
	}
}

func requireConfirmation(
	ctx tool.Context, t tool.Tool, args map[string]any,
) (map[string]any, error) {
	if t.Name() == agentDefs["Aagje"].name {
		return nil, nil
	}

	confirmation := ctx.ToolConfirmation()
	if confirmation == nil {
		hint := fmt.Sprintf(
			"The agent wants to call tool %q with args %v. Do you approve?",
			t.Name(), args,
		)
		if err := ctx.RequestConfirmation(hint, nil); err != nil {
			return nil, fmt.Errorf("failed to request confirmation: %w", err)
		}
		return nil, nil
	}

	if !confirmation.Confirmed {
		return map[string]any{
			"error": fmt.Sprintf("tool %q was denied by the user", t.Name()),
		}, nil
	}
	return nil, nil
}

func buildOrchestrator(m model.LLM, mcpToolSet tool.Toolset) (
	agent.Agent, error,
) {
	def := agentDefs["Aagje"]
	retrieval, err := llmagent.New(
		llmagent.Config{
			Name:        def.name,
			Model:       m,
			Description: def.description,
			Instruction: *def.instruction,
			Tools:       []tool.Tool{geminitool.GoogleSearch{}},
			Toolsets: []tool.Toolset{
				tool.FilterToolset(mcpToolSet, tool.StringPredicate(def.toolNames)),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s agent: %w", def.name, err)
	}

	retrievalTool := agenttool.New(retrieval, nil)
	agents := make(map[string]agent.Agent, 2)
	for _, key := range []string{"Knorretje", "Daan"} {
		def = agentDefs[key]
		a, err := llmagent.New(
			llmagent.Config{
				Name:        def.name,
				Model:       m,
				Description: def.description,
				Instruction: *def.instruction,
				Tools:       []tool.Tool{retrievalTool},
				Toolsets: []tool.Toolset{
					tool.FilterToolset(mcpToolSet, tool.StringPredicate(def.toolNames)),
				},
				BeforeToolCallbacks: []llmagent.BeforeToolCallback{
					requireConfirmation,
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s agent: %w", def.name, err)
		}
		agents[key] = a
	}

	oDef := agentDefs["Nina"]
	orchestrator, err := llmagent.New(
		llmagent.Config{
			Name:        oDef.name,
			Model:       m,
			Description: oDef.description,
			Instruction: *oDef.instruction,
			Tools:       []tool.Tool{retrievalTool},
			SubAgents:   []agent.Agent{agents["Knorretje"], agents["Daan"]},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s agent: %w", oDef.name, err)
	}

	return orchestrator, nil
}

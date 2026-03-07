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
		description: "Orchestrates YouTube workflows by planning multi-step tasks and delegating to specialized agents.",
		instruction: &orchestratorInstruction,
		envKey:      "YUTU_AGENT_INSTRUCTION",
	},
	"Aagje": {
		name:        "Retrieval",
		description: "Retrieves and analyzes information from YouTube — listing videos, channels, playlists, comments, analytics, and more. Use for read-only queries, research, and analysis.",
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
		description: "Creates and updates YouTube content — posting comments, uploading videos, creating playlists, updating metadata, setting thumbnails. Can also look up channels, videos, comments, playlists, captions, and subscriptions to gather context needed for modifications.",
		instruction: &modifierInstruction,
		envKey:      "YUTU_MODIFIER_INSTRUCTION",
		toolNames: []string{
			"caption-list",
			"channel-list",
			"comment-list",
			"commentThread-list",
			"playlist-list",
			"playlistItem-list",
			"search-list",
			"subscription-list",
			"video-list",
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
		description: "Deletes YouTube content — videos, playlists, comments, captions, subscriptions. Can search and verify targets before deletion. Handles destructive operations requiring extra caution.",
		instruction: &destroyerInstruction,
		envKey:      "YUTU_DESTROYER_INSTRUCTION",
		toolNames: []string{
			"caption-list",
			"channel-list",
			"comment-list",
			"commentThread-list",
			"playlist-list",
			"playlistItem-list",
			"search-list",
			"subscription-list",
			"video-list",
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

var confirmationToolNames = computeConfirmationToolNames()

func computeConfirmationToolNames() map[string]struct{} {
	readOnly := make(map[string]struct{}, len(agentDefs["Aagje"].toolNames))
	for _, name := range agentDefs["Aagje"].toolNames {
		readOnly[name] = struct{}{}
	}

	result := make(map[string]struct{})
	for _, key := range []string{"Knorretje", "Daan"} {
		for _, name := range agentDefs[key].toolNames {
			if _, ok := readOnly[name]; !ok {
				result[name] = struct{}{}
			}
		}
	}
	return result
}

func init() {
	for _, def := range agentDefs {
		if v, ok := os.LookupEnv(def.envKey); ok && v != "" {
			*def.instruction = v
		}
	}
}

func requireConfirmation(name string, _ any) bool {
	_, ok := confirmationToolNames[name]
	return ok
}

func buildOrchestrator(
	advancedModel, liteModel model.LLM, mcpToolSet tool.Toolset,
) (
	agent.Agent, error,
) {
	subAgents := make([]agent.Agent, 0, 3)
	for _, key := range []string{"Aagje", "Knorretje", "Daan"} {
		def := agentDefs[key]
		var extraTools []tool.Tool
		if key == "Aagje" {
			extraTools = []tool.Tool{geminitool.GoogleSearch{}}
		}
		a, err := llmagent.New(
			llmagent.Config{
				Name:        def.name,
				Model:       liteModel,
				Description: def.description,
				Instruction: *def.instruction,
				Tools:       extraTools,
				Toolsets: []tool.Toolset{
					tool.FilterToolset(mcpToolSet, tool.StringPredicate(def.toolNames)),
				},
				DisallowTransferToParent: true,
				DisallowTransferToPeers:  true,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s agent: %w", def.name, err)
		}
		subAgents = append(subAgents, a)
	}

	oDef := agentDefs["Nina"]
	orchestrator, err := llmagent.New(
		llmagent.Config{
			Name:        oDef.name,
			Model:       advancedModel,
			Description: oDef.description,
			Instruction: *oDef.instruction,
			Tools:       []tool.Tool{geminitool.GoogleSearch{}},
			SubAgents:   subAgents,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s agent: %w", oDef.name, err)
	}

	return orchestrator, nil
}

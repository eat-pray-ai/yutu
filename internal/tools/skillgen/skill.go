// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"
)

const skillDescription = "Use whenever the user mentions YouTube, video uploads, channel management, playlists, video SEO, or any YouTube Data API operation. Manages videos, playlists, comments, captions, subscriptions, thumbnails, analytics, and more."

const growthTips = `- **Titles**: Curiosity gaps + power words. Front-load keywords. Under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Consistent schedule matters.
- **Engagement**: Pin a comment with a question. Reply within the first hour.`

const safetyRules = `- **Confirm before deleting**: Always list exactly what will be deleted and get explicit user confirmation before calling any destructive tool.
- **Verify targets**: Before modifying or deleting, use a read tool to confirm the resource exists and matches the user's intent.
- **Preserve existing data**: When updating resources, only modify fields the user explicitly asked to change.
- **Report resource IDs**: Always include resource IDs in responses so the user can reference them.`

const workflow = `For every user request:

1. **Classify**: Is this a read-only query, a content modification, a deletion, or a multistep workflow?
2. **Gather context**: Use read tools to find resource IDs, verify targets, and understand current state.
3. **Execute**: Use write or delete tools with the gathered data.
4. **Report**: State what was done, including resource IDs and titles.`

type toolClass int

const (
	classRead toolClass = iota
	classWrite
	classDelete
)

var destructiveVerbs = map[string]bool{
	"delete": true, "unset": true,
}

var readOnlyVerbs = map[string]bool{
	"list": true, "download": true, "getRating": true,
}

func classifyVerb(verb string) toolClass {
	if destructiveVerbs[verb] {
		return classDelete
	}
	if readOnlyVerbs[verb] {
		return classRead
	}
	return classWrite
}

func writeFrontmatter(b *strings.Builder, t skillTarget) {
	_, _ = fmt.Fprintf(
		b, "---\nname: youtube\ndescription: \"%s\"\nlicense: MIT\n",
		skillDescription,
	)
	if t.compatibility != "" {
		_, _ = fmt.Fprintf(b, "compatibility: %s\n", t.compatibility)
	}
	b.WriteString("metadata:\n  author: eat-pray-ai\n  homepage: \"https://github.com/eat-pray-ai/yutu\"\n---\n\n")
}

func writeSkill(path string, resources []resourceEntry, t skillTarget) error {
	var b strings.Builder

	writeFrontmatter(&b, t)

	b.WriteString("# YouTube\n\n")
	names := make([]string, len(resources))
	for i, r := range resources {
		names[i] = r.name
	}
	_, _ = fmt.Fprintf(&b, "Manage YouTube resources via MCP tool calls — %s.\n\n", strings.Join(names, ", "))

	if t.includeSetup {
		b.WriteString("## Setup\n\n")
		b.WriteString("Ensure `yutu` is installed and running as an MCP server. See [references/setup.md](references/setup.md).\n\n")
	}

	b.WriteString("## Key Principles\n\n")
	b.WriteString(safetyRules + "\n\n")

	var readTools, writeTools, deleteTools []string
	for _, r := range resources {
		for _, v := range r.verbs {
			toolName := r.name + "-" + v.name
			switch classifyVerb(v.name) {
			case classRead:
				readTools = append(readTools, toolName)
			case classWrite:
				writeTools = append(writeTools, toolName)
			case classDelete:
				deleteTools = append(deleteTools, toolName)
			}
		}
	}

	b.WriteString("## Tools\n\n")
	b.WriteString("### Read (safe, no side effects)\n\n")
	for _, t := range readTools {
		_, _ = fmt.Fprintf(&b, "- `%s`\n", t)
	}

	b.WriteString("\n### Write (create or update content)\n\n")
	for _, t := range writeTools {
		_, _ = fmt.Fprintf(&b, "- `%s`\n", t)
	}

	b.WriteString("\n### Destructive (irreversible deletions — confirm with user first)\n\n")
	for _, t := range deleteTools {
		_, _ = fmt.Fprintf(&b, "- `%s`\n", t)
	}

	b.WriteString("\n## Common Workflows\n\n")
	b.WriteString("See [references/workflows.md](references/workflows.md) for step-by-step walkthroughs.\n\n")

	b.WriteString("## YouTube Growth Tips\n\n")
	b.WriteString("See [references/seo-guide.md](references/seo-guide.md) for the full guide. When creating or updating content, apply these principles:\n\n")
	b.WriteString(growthTips + "\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeInstruction(path string) error {
	var b strings.Builder

	b.WriteString("You are an expert YouTube growth strategist and workflow assistant. You are authenticated to operate YouTube channels on behalf of users via tool calls.\n\n")
	b.WriteString("Consult the **youtube** skill for the full tool catalog, safety rules, workflows, and SEO guidelines.\n\n")
	b.WriteString("## Workflow\n\n")
	b.WriteString(workflow + "\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

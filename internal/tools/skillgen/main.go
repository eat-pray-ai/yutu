// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"

	// Blank-import every resource package so that its init() registers the
	// subcommand (and sub-subcommands) on cmd.RootCmd.
	_ "github.com/eat-pray-ai/yutu/cmd/activity"
	_ "github.com/eat-pray-ai/yutu/cmd/caption"
	_ "github.com/eat-pray-ai/yutu/cmd/channel"
	_ "github.com/eat-pray-ai/yutu/cmd/channelBanner"
	_ "github.com/eat-pray-ai/yutu/cmd/channelSection"
	_ "github.com/eat-pray-ai/yutu/cmd/comment"
	_ "github.com/eat-pray-ai/yutu/cmd/commentThread"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nLanguage"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nRegion"
	_ "github.com/eat-pray-ai/yutu/cmd/liveChatBan"
	_ "github.com/eat-pray-ai/yutu/cmd/member"
	_ "github.com/eat-pray-ai/yutu/cmd/membershipsLevel"
	_ "github.com/eat-pray-ai/yutu/cmd/playlist"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistImage"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistItem"
	_ "github.com/eat-pray-ai/yutu/cmd/search"
	_ "github.com/eat-pray-ai/yutu/cmd/subscription"
	_ "github.com/eat-pray-ai/yutu/cmd/superChatEvent"
	_ "github.com/eat-pray-ai/yutu/cmd/thirdPartyLink"
	_ "github.com/eat-pray-ai/yutu/cmd/thumbnail"
	_ "github.com/eat-pray-ai/yutu/cmd/video"
	_ "github.com/eat-pray-ai/yutu/cmd/videoAbuseReportReason"
	_ "github.com/eat-pray-ai/yutu/cmd/videoCategory"
	_ "github.com/eat-pray-ai/yutu/cmd/watermark"
)

//go:embed setup.md
var setupContent string

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

type verbEntry struct {
	name  string
	short string
}

type resourceEntry struct {
	name  string
	kebab string
	human string
	short string
	long  string
	verbs []verbEntry
}

type categoryGroup struct {
	name      string
	resources []resourceEntry
}

// ---------------------------------------------------------------------------
// Category mapping
// ---------------------------------------------------------------------------

var categoryOrder = []string{
	"Content", "Organization", "Community", "Channel", "Discovery", "Metadata",
}

var categoryMap = map[string]string{
	"video":                  "Content",
	"caption":                "Content",
	"thumbnail":              "Content",
	"watermark":              "Content",
	"playlist":               "Organization",
	"playlistItem":           "Organization",
	"playlistImage":          "Organization",
	"comment":                "Community",
	"commentThread":          "Community",
	"subscription":           "Community",
	"member":                 "Community",
	"membershipsLevel":       "Community",
	"liveChatBan":            "Community",
	"superChatEvent":         "Community",
	"channel":                "Channel",
	"channelBanner":          "Channel",
	"channelSection":         "Channel",
	"thirdPartyLink":         "Channel",
	"search":                 "Discovery",
	"activity":               "Discovery",
	"videoCategory":          "Metadata",
	"videoAbuseReportReason": "Metadata",
	"i18nLanguage":           "Metadata",
	"i18nRegion":             "Metadata",
}

// resourceCategory returns the category for a resource, or "Other" if unmapped.
func resourceCategory(name string) string {
	if cat, ok := categoryMap[name]; ok {
		return cat
	}
	return "Other"
}

// collectResource builds a resourceEntry from a cobra resource command.
func collectResource(c *cobra.Command) resourceEntry {
	name := c.Name()
	var verbs []verbEntry
	for _, sub := range c.Commands() {
		if sub.Name() == "help" {
			continue
		}
		verbs = append(verbs, verbEntry{name: sub.Name(), short: sub.Short})
	}
	sort.Slice(
		verbs, func(i, j int) bool {
			return verbs[i].name < verbs[j].name
		},
	)
	return resourceEntry{
		name:  name,
		kebab: camelToKebab(name),
		human: camelToWords(name),
		short: c.Short,
		long:  c.Long,
		verbs: verbs,
	}
}

// groupByCategory groups resources into categories ordered by categoryOrder,
// with "Other" appended at the end for any unmapped resources.
func groupByCategory(resources []resourceEntry) []categoryGroup {
	byCategory := make(map[string][]resourceEntry)
	for _, r := range resources {
		cat := resourceCategory(r.name)
		byCategory[cat] = append(byCategory[cat], r)
	}

	var groups []categoryGroup
	for _, cat := range categoryOrder {
		if rs, ok := byCategory[cat]; ok {
			sort.Slice(
				rs, func(i, j int) bool {
					return rs[i].name < rs[j].name
				},
			)
			groups = append(groups, categoryGroup{name: cat, resources: rs})
			delete(byCategory, cat)
		}
	}
	// Append "Other" for any remaining unmapped categories.
	if rs, ok := byCategory["Other"]; ok {
		sort.Slice(
			rs, func(i, j int) bool {
				return rs[i].name < rs[j].name
			},
		)
		groups = append(groups, categoryGroup{name: "Other", resources: rs})
	}
	return groups
}

// ---------------------------------------------------------------------------
// Unified description builder
// ---------------------------------------------------------------------------

// buildUnifiedDescription constructs the skill description with high-frequency
// triggers and a broad fallback. Intentionally concise — no exhaustive verb list.
func buildUnifiedDescription(_ []categoryGroup) string {
	return "Use whenever the user mentions YouTube, video uploads, channel management, playlists, video SEO, or any YouTube Data API operation. Manages videos, playlists, comments, captions, subscriptions, thumbnails, analytics, and more via the yutu CLI."
}

// ---------------------------------------------------------------------------
// Static content for SKILL.md sections
// ---------------------------------------------------------------------------

const workflowSummary = `| Task | Quick Command |
|------|---------------|
| Upload a video | ` + "`yutu video insert --file video.mp4 --title \"...\" --privacy public`" + ` |
| Update video metadata | ` + "`yutu video list --ids ID`" + ` then ` + "`yutu video update --id ID --title \"...\"`" + ` |
| Create playlist + add videos | ` + "`yutu playlist insert`" + ` → ` + "`yutu playlistItem insert`" + ` |
| Post a comment | ` + "`yutu commentThread insert --channelId ... --videoId ... --textOriginal \"...\"`" + ` |
| Channel analytics | ` + "`yutu channel list --mine --output json`" + ` |
| Competitor analysis | ` + "`yutu channel list --forHandle @handle --output json`" + ` |
| Delete content | Always ` + "`list`" + ` first, then ` + "`delete`" + ` — irreversible |
| Subscribe/unsubscribe | Check ` + "`yutu subscription list --mine --forChannelId ...`" + ` before acting |`

const growthTips = `- **Titles**: Curiosity gaps + power words. Front-load keywords. Under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Consistent schedule matters.
- **Engagement**: Pin a comment with a question. Reply within the first hour.`

// ---------------------------------------------------------------------------
// Unified SKILL.md writer
// ---------------------------------------------------------------------------

// writeUnifiedSkill generates the single unified SKILL.md for the youtube skill.
func writeUnifiedSkill(path string, groups []categoryGroup) error {
	var b strings.Builder

	desc := buildUnifiedDescription(groups)

	_, _ = fmt.Fprintf(
		&b, `---
name: youtube
description: "%s"
license: MIT
compatibility: Requires the yutu CLI binary (installable via npm, brew, or winget) and Google Cloud OAuth credentials for YouTube Data API v3.
metadata:
  author: eat-pray-ai
  homepage: "https://github.com/eat-pray-ai/yutu"
---

`, strings.ReplaceAll(desc, `"`, `\"`),
	)

	b.WriteString("# YouTube\n\n")
	b.WriteString("Manage YouTube resources using the `yutu` CLI — videos, playlists, comments, channels, captions, subscriptions, and more.\n\n")

	b.WriteString("## Quick Start\n\n")
	b.WriteString("1. Ensure `yutu` is installed and authenticated. If not, follow [references/setup.md](references/setup.md).\n")
	b.WriteString("2. Identify the resource and operation from the tables below.\n")
	b.WriteString("3. Run `yutu <resource> <operation> -h` for full flag details on any command.\n")
	b.WriteString("4. For multistep tasks (upload + thumbnail + playlist), see [references/workflows.md](references/workflows.md).\n\n")

	b.WriteString("## Key Principles\n\n")
	b.WriteString("- Always verify before destructive operations — deletions are irreversible.\n")
	b.WriteString("- Use `--output json` when you need to parse or chain results.\n")
	b.WriteString("- Get your channel ID with `yutu channel list --mine` — many operations need it.\n")
	b.WriteString("- When updating metadata, only specify the fields you want to change.\n\n")

	b.WriteString("## Operations\n\n")

	for _, g := range groups {
		_, _ = fmt.Fprintf(&b, "### %s\n\n", g.name)
		b.WriteString("| Resource | Operation | Description |\n")
		b.WriteString("|----------|-----------|-------------|\n")
		for _, r := range g.resources {
			for _, v := range r.verbs {
				_, _ = fmt.Fprintf(
					&b, "| %s | %s | %s |\n", r.human, v.name, escPipe(v.short),
				)
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("## Common Workflows\n\n")
	b.WriteString("See [references/workflows.md](references/workflows.md) for step-by-step walkthroughs of each task below.\n\n")
	b.WriteString(workflowSummary + "\n\n")

	b.WriteString("## YouTube Growth Tips\n\n")
	b.WriteString("See [references/seo-guide.md](references/seo-guide.md) for the full guide. When uploading or updating video metadata, apply these principles:\n\n")
	b.WriteString(growthTips + "\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// ---------------------------------------------------------------------------
// Setup writer
// ---------------------------------------------------------------------------

func writeSetup(path string) error {
	return os.WriteFile(path, []byte(setupContent), 0o644)
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	out := flag.String("out", "./skills", "output directory for generated skills")
	flag.Parse()

	root := cmd.RootCmd
	root.InitDefaultHelpCmd()

	dir := filepath.Join(*out, "youtube")
	refDir := filepath.Join(dir, "references")

	if err := os.MkdirAll(refDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", refDir, err)
	}

	setupPath := filepath.Join(refDir, "setup.md")
	if err := writeSetup(setupPath); err != nil {
		log.Fatalf("write setup %s: %v", setupPath, err)
	}

	var resources []resourceEntry
	for _, c := range root.Commands() {
		if !strings.HasPrefix(c.Short, "Manage") {
			continue
		}
		resources = append(resources, collectResource(c))
	}

	groups := groupByCategory(resources)

	skillPath := filepath.Join(dir, "SKILL.md")
	if err := writeUnifiedSkill(skillPath, groups); err != nil {
		log.Fatalf("write skill %s: %v", skillPath, err)
	}

	totalVerbs := 0
	for _, g := range groups {
		for _, r := range g.resources {
			totalVerbs += len(r.verbs)
		}
	}
	fmt.Printf(
		"Generated unified youtube skill: %d resources, %d verbs, %d categories\n",
		len(resources), totalVerbs, len(groups),
	)
}

// ---------------------------------------------------------------------------
// String utilities
// ---------------------------------------------------------------------------

// escPipe escapes pipe characters for markdown table cells.
func escPipe(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

// camelToWords splits a camelCase string into lowercase space-separated words.
func camelToWords(s string) string {
	return strings.Join(camelSplit(s), " ")
}

func camelToKebab(s string) string {
	return strings.Join(camelSplit(s), "-")
}

func camelSplit(s string) []string {
	var words []string
	var cur []byte
	for i := range len(s) {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' && len(cur) > 0 {
			words = append(words, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
		cur = append(cur, ch)
	}
	if len(cur) > 0 {
		words = append(words, strings.ToLower(string(cur)))
	}
	return words
}

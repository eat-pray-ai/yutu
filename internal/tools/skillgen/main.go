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
	_ "github.com/eat-pray-ai/yutu/cmd/member"
	_ "github.com/eat-pray-ai/yutu/cmd/membershipsLevel"
	_ "github.com/eat-pray-ai/yutu/cmd/playlist"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistImage"
	_ "github.com/eat-pray-ai/yutu/cmd/playlistItem"
	_ "github.com/eat-pray-ai/yutu/cmd/search"
	_ "github.com/eat-pray-ai/yutu/cmd/subscription"
	_ "github.com/eat-pray-ai/yutu/cmd/superChatEvent"
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
	"superChatEvent":         "Community",
	"channel":                "Channel",
	"channelBanner":          "Channel",
	"channelSection":         "Channel",
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
	sort.Slice(verbs, func(i, j int) bool {
		return verbs[i].name < verbs[j].name
	})
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
			sort.Slice(rs, func(i, j int) bool {
				return rs[i].name < rs[j].name
			})
			groups = append(groups, categoryGroup{name: cat, resources: rs})
			delete(byCategory, cat)
		}
	}
	// Append "Other" for any remaining unmapped categories.
	if rs, ok := byCategory["Other"]; ok {
		sort.Slice(rs, func(i, j int) bool {
			return rs[i].name < rs[j].name
		})
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
	return "Use when working with YouTube — upload videos, search content, manage playlists and channels, post and moderate comments, handle subscriptions and memberships, add captions, set thumbnails, check analytics, or any YouTube Data API operation via the yutu CLI."
}

// ---------------------------------------------------------------------------
// Static content for SKILL.md sections
// ---------------------------------------------------------------------------

const workflowSummary = `- **Upload a video**: ` + "`yutu video insert --file video.mp4 --title \"...\" --privacy public`" + `, then optionally set thumbnail
- **Update video metadata**: Fetch current with ` + "`yutu video list --id VIDEO_ID`" + `, then update changed fields
- **Create playlist + add videos**: Create with ` + "`yutu playlist insert`" + `, find videos with ` + "`yutu search list --forMine`" + `, add with ` + "`yutu playlistItem insert`" + `
- **Post a comment**: Get channel ID with ` + "`yutu channel list --mine`" + `, find video, then ` + "`yutu commentThread insert`" + `
- **Channel analytics**: ` + "`yutu channel list --mine`" + ` + ` + "`yutu search list --forMine`" + ` + ` + "`yutu video list --id ...`" + `
- **Competitor analysis**: ` + "`yutu channel list --forHandle @handle`" + ` + compare stats and top videos
- **Delete content**: Always verify with a list command first, then delete — deletions are irreversible
- **Subscribe/unsubscribe**: Check with ` + "`yutu subscription list --mine --forChannelId ...`" + ` before acting`

const growthTips = `- **Titles**: Use curiosity gaps and power words. Front-load keywords. Keep under 60 characters.
- **Descriptions**: First 2 lines appear in search. Include keywords, timestamps, CTAs, and 3-5 hashtags.
- **Tags**: Mix broad and long-tail keywords. First 2-3 tags carry the most weight.
- **Thumbnails**: High contrast, 3-4 word text, expressive faces, consistent branding.
- **Publishing**: Post when audience is active. Maintain consistent schedule.
- **Engagement**: Pin a comment with a question. Reply within the first hour.`

// ---------------------------------------------------------------------------
// Unified SKILL.md writer
// ---------------------------------------------------------------------------

// writeUnifiedSkill generates the single unified SKILL.md for the youtube skill.
func writeUnifiedSkill(path string, groups []categoryGroup) error {
	var b strings.Builder

	desc := buildUnifiedDescription(groups)

	_, _ = fmt.Fprintf(&b, `---
name: youtube
description: "%s"
metadata:
  openclaw:
    requires:
      env:
        - YUTU_CREDENTIAL
        - YUTU_CACHE_TOKEN
      bins:
        - yutu
      config:
        - client_secret.json
        - youtube.token.json
    primaryEnv: YUTU_CREDENTIAL
    emoji: "\U0001F3AC\U0001F430"
    homepage: https://github.com/eat-pray-ai/yutu
    install:
      - kind: node
        package: "@eat-pray-ai/yutu"
        bins: [yutu]
---

`, strings.ReplaceAll(desc, `"`, `\"`))

	b.WriteString("# YouTube\n\n")
	b.WriteString("Manage YouTube resources using the yutu CLI — videos, playlists, comments, channels, captions, subscriptions, and more.\n\n")

	b.WriteString("## Before You Begin\n\n")
	b.WriteString("yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. ")
	b.WriteString("If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.\n\n")

	b.WriteString("## Operations\n\n")
	b.WriteString("Run `yutu <resource> <verb> -h` for full flag details and examples.\n\n")

	for _, g := range groups {
		_, _ = fmt.Fprintf(&b, "### %s\n\n", g.name)
		b.WriteString("| Resource | Operation | Description |\n")
		b.WriteString("|----------|-----------|-------------|\n")
		for _, r := range g.resources {
			for _, v := range r.verbs {
				_, _ = fmt.Fprintf(&b, "| %s | %s | %s |\n", r.human, v.name, escPipe(v.short))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("## Common Workflows\n\n")
	b.WriteString("See [references/workflows.md](references/workflows.md) for detailed walkthroughs.\n\n")
	b.WriteString(workflowSummary + "\n\n")

	b.WriteString("## YouTube Growth Tips\n\n")
	b.WriteString("See [references/seo-guide.md](references/seo-guide.md) for the full guide.\n\n")
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
	fmt.Printf("Generated unified youtube skill: %d resources, %d verbs, %d categories\n",
		len(resources), totalVerbs, len(groups))
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

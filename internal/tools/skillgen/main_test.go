// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResourceCategory(t *testing.T) {
	expected := map[string]string{
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

	for resource, wantCategory := range expected {
		got := resourceCategory(resource)
		if got != wantCategory {
			t.Errorf("resourceCategory(%q) = %q, want %q", resource, got, wantCategory)
		}
	}

	// Unknown resource should return "Other".
	if got := resourceCategory("unknownResource"); got != "Other" {
		t.Errorf("resourceCategory(%q) = %q, want %q", "unknownResource", got, "Other")
	}
}

func TestGroupByCategory(t *testing.T) {
	resources := []resourceEntry{
		{name: "video", verbs: []verbEntry{{name: "list", short: "List videos"}}},
		{name: "playlist", verbs: []verbEntry{{name: "insert", short: "Create a playlist"}}},
		{name: "comment", verbs: []verbEntry{{name: "list", short: "List comments"}}},
		{name: "channel", verbs: []verbEntry{{name: "list", short: "List channels"}}},
		{name: "search", verbs: []verbEntry{{name: "list", short: "Search"}}},
		{name: "videoCategory", verbs: []verbEntry{{name: "list", short: "List categories"}}},
		{name: "caption", verbs: []verbEntry{{name: "list", short: "List captions"}}},
	}

	groups := groupByCategory(resources)

	// Should follow categoryOrder: Content, Organization, Community, Channel, Discovery, Metadata.
	wantOrder := []string{"Content", "Organization", "Community", "Channel", "Discovery", "Metadata"}
	if len(groups) != len(wantOrder) {
		t.Fatalf("got %d groups, want %d", len(groups), len(wantOrder))
	}
	for i, g := range groups {
		if g.name != wantOrder[i] {
			t.Errorf("group[%d].name = %q, want %q", i, g.name, wantOrder[i])
		}
	}

	// Content should have video and caption (sorted).
	contentGroup := groups[0]
	if len(contentGroup.resources) != 2 {
		t.Fatalf("Content group has %d resources, want 2", len(contentGroup.resources))
	}
	if contentGroup.resources[0].name != "caption" {
		t.Errorf("Content[0] = %q, want %q", contentGroup.resources[0].name, "caption")
	}
	if contentGroup.resources[1].name != "video" {
		t.Errorf("Content[1] = %q, want %q", contentGroup.resources[1].name, "video")
	}
}

func TestGroupByCategoryWithOther(t *testing.T) {
	resources := []resourceEntry{
		{name: "video", verbs: []verbEntry{{name: "list", short: "List videos"}}},
		{name: "unknownThing", verbs: []verbEntry{{name: "list", short: "List unknowns"}}},
	}

	groups := groupByCategory(resources)

	// Should have Content + Other.
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}
	if groups[0].name != "Content" {
		t.Errorf("groups[0].name = %q, want %q", groups[0].name, "Content")
	}
	if groups[1].name != "Other" {
		t.Errorf("groups[1].name = %q, want %q", groups[1].name, "Other")
	}
	if groups[1].resources[0].name != "unknownThing" {
		t.Errorf("Other[0] = %q, want %q", groups[1].resources[0].name, "unknownThing")
	}
}

func TestBuildUnifiedDescription(t *testing.T) {
	groups := []categoryGroup{
		{
			name: "Content",
			resources: []resourceEntry{
				{name: "video", human: "video", verbs: []verbEntry{
					{name: "list", short: "List videos"},
					{name: "insert", short: "Upload a video"},
				}},
			},
		},
		{
			name: "Community",
			resources: []resourceEntry{
				{name: "comment", human: "comment", verbs: []verbEntry{
					{name: "list", short: "List comments"},
				}},
			},
		},
	}

	desc := buildUnifiedDescription(groups)

	// Should start with the broad fallback.
	if !strings.HasPrefix(desc, "Use when working with YouTube") {
		t.Error("description should start with broad fallback")
	}

	// Should contain key topics.
	for _, phrase := range []string{
		"upload videos",
		"search content",
		"manage playlists",
		"captions",
		"thumbnails",
		"analytics",
		"YouTube Data API",
	} {
		if !strings.Contains(desc, phrase) {
			t.Errorf("description missing phrase %q", phrase)
		}
	}

	// Should be concise — no duplicated trigger list.
	if strings.Contains(desc, "Triggers:") {
		t.Error("description should not have a separate Triggers section")
	}
}

func TestWriteUnifiedSkill(t *testing.T) {
	groups := []categoryGroup{
		{
			name: "Content",
			resources: []resourceEntry{
				{name: "video", human: "video", verbs: []verbEntry{
					{name: "insert", short: "Upload a video"},
					{name: "list", short: "List videos"},
				}},
				{name: "caption", human: "caption", verbs: []verbEntry{
					{name: "list", short: "List captions"},
				}},
			},
		},
		{
			name: "Organization",
			resources: []resourceEntry{
				{name: "playlist", human: "playlist", verbs: []verbEntry{
					{name: "insert", short: "Create a playlist"},
				}},
			},
		},
		{
			name: "Community",
			resources: []resourceEntry{
				{name: "comment", human: "comment", verbs: []verbEntry{
					{name: "list", short: "List comments"},
				}},
			},
		},
		{
			name: "Channel",
			resources: []resourceEntry{
				{name: "channel", human: "channel", verbs: []verbEntry{
					{name: "list", short: "List channels"},
				}},
			},
		},
		{
			name: "Discovery",
			resources: []resourceEntry{
				{name: "search", human: "search", verbs: []verbEntry{
					{name: "list", short: "Search YouTube"},
				}},
			},
		},
		{
			name: "Metadata",
			resources: []resourceEntry{
				{name: "videoCategory", human: "video category", verbs: []verbEntry{
					{name: "list", short: "List video categories"},
				}},
			},
		},
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "SKILL.md")

	if err := writeUnifiedSkill(path, groups); err != nil {
		t.Fatalf("writeUnifiedSkill: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}
	content := string(data)

	// Frontmatter has name: youtube.
	if !strings.Contains(content, "name: youtube") {
		t.Error("SKILL.md missing 'name: youtube' in frontmatter")
	}

	// All 6 category headers present.
	for _, cat := range []string{"Content", "Organization", "Community", "Channel", "Discovery", "Metadata"} {
		if !strings.Contains(content, "### "+cat) {
			t.Errorf("SKILL.md missing category header ### %s", cat)
		}
	}

	// Operations table rows present.
	for _, row := range []string{
		"| video | insert | Upload a video |",
		"| video | list | List videos |",
		"| caption | list | List captions |",
		"| playlist | insert | Create a playlist |",
		"| comment | list | List comments |",
		"| channel | list | List channels |",
		"| search | list | Search YouTube |",
		"| video category | list | List video categories |",
	} {
		if !strings.Contains(content, row) {
			t.Errorf("SKILL.md missing table row %q", row)
		}
	}

	// CLI help instruction present.
	if !strings.Contains(content, "Run `yutu <resource> <verb> -h` for full flag details and examples.") {
		t.Error("SKILL.md missing CLI help instruction")
	}

	// Workflow section present.
	if !strings.Contains(content, "## Common Workflows") {
		t.Error("SKILL.md missing Common Workflows section")
	}
	if !strings.Contains(content, "references/workflows.md") {
		t.Error("SKILL.md missing workflows.md reference")
	}
	if !strings.Contains(content, "Upload a video") {
		t.Error("SKILL.md missing workflow summary content")
	}

	// SEO/Growth Tips section present.
	if !strings.Contains(content, "## YouTube Growth Tips") {
		t.Error("SKILL.md missing YouTube Growth Tips section")
	}
	if !strings.Contains(content, "references/seo-guide.md") {
		t.Error("SKILL.md missing seo-guide.md reference")
	}
	if !strings.Contains(content, "curiosity gaps") {
		t.Error("SKILL.md missing growth tips content")
	}
}

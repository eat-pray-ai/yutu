// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	got := generate()

	checks := []struct {
		name    string
		substr  string
	}{
		{"shebang", "#!/usr/bin/env bash"},
		{"set strict", "set -euo pipefail"},
		{"root help", `"$YUTU_PATH" help`},
		{"root completion", `"$YUTU_PATH" completion`},
		{"root version", `"$YUTU_PATH" version`},
		{"auth section", `echo "======= auth ======="`},
		{"mcp section", `echo "======= mcp ======="`},
		{"agent section", `echo "======= agent ======="`},
		{"video section", `echo "======= video ======="`},
		{"video list", `"$YUTU_PATH" video list --help`},
		{"video insert", `"$YUTU_PATH" video insert --help`},
		{"video delete", `"$YUTU_PATH" video delete --help`},
		{"channel section", `echo "======= channel ======="`},
		{"playlist section", `echo "======= playlist ======="`},
		{"comment section", `echo "======= comment ======="`},
		{"caption section", `echo "======= caption ======="`},
		{"search section", `echo "======= search ======="`},
	}

	for _, c := range checks {
		if !strings.Contains(got, c.substr) {
			t.Errorf("%s: expected output to contain %q", c.name, c.substr)
		}
	}

	if !strings.HasSuffix(got, "\n") {
		t.Error("expected output to end with newline")
	}
}

func TestGenerateNoResourceMissed(t *testing.T) {
	got := generate()

	resources := []string{
		"activity", "caption", "channel", "channelBanner", "channelSection",
		"comment", "commentThread", "i18nLanguage", "i18nRegion", "member",
		"membershipsLevel", "playlist", "playlistImage", "playlistItem",
		"search", "subscription", "superChatEvent", "thumbnail", "video",
		"videoAbuseReportReason", "videoCategory", "watermark",
	}

	for _, r := range resources {
		section := `echo "======= ` + r + ` ======="`
		if !strings.Contains(got, section) {
			t.Errorf("missing resource section for %s", r)
		}
		helpLine := `"$YUTU_PATH" ` + r + ` --help`
		if !strings.Contains(got, helpLine) {
			t.Errorf("missing --help for resource %s", r)
		}
	}
}

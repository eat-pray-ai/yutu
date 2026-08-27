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

	"github.com/eat-pray-ai/yutu/cmd"

	// Blank-import every resource package so that its init() registers the
	// subcommand (and sub-subcommands) on cmd.RootCmd.
	_ "github.com/eat-pray-ai/yutu/cmd/abuseReport"
	_ "github.com/eat-pray-ai/yutu/cmd/activity"
	_ "github.com/eat-pray-ai/yutu/cmd/caption"
	_ "github.com/eat-pray-ai/yutu/cmd/channel"
	_ "github.com/eat-pray-ai/yutu/cmd/channelBanner"
	_ "github.com/eat-pray-ai/yutu/cmd/channelSection"
	_ "github.com/eat-pray-ai/yutu/cmd/comment"
	_ "github.com/eat-pray-ai/yutu/cmd/commentThread"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nLanguage"
	_ "github.com/eat-pray-ai/yutu/cmd/i18nRegion"
	_ "github.com/eat-pray-ai/yutu/cmd/liveBroadcast"
	_ "github.com/eat-pray-ai/yutu/cmd/liveChatBan"
	_ "github.com/eat-pray-ai/yutu/cmd/liveChatMessage"
	_ "github.com/eat-pray-ai/yutu/cmd/liveChatModerator"
	_ "github.com/eat-pray-ai/yutu/cmd/liveStream"
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

//go:embed workflows.md
var workflowsContent string

//go:embed seo-guide.md
var seoGuideContent string

type skillTarget struct {
	subdir        string
	refs          map[string]string
	compatibility string
	includeSetup  bool
}

var extSkill = skillTarget{
	subdir: filepath.Join("skills", "youtube"),
	refs: map[string]string{
		"setup.md":     setupContent,
		"workflows.md": workflowsContent,
		"seo-guide.md": seoGuideContent,
	},
	compatibility: "Requires the yutu CLI binary (installable via npm, brew, or winget) and Google Cloud OAuth credentials for YouTube Data API v3.",
	includeSetup:  true,
}

var agentSkill = skillTarget{
	subdir: filepath.Join("skills", "youtube"),
	refs: map[string]string{
		"workflows.md": workflowsContent,
		"seo-guide.md": seoGuideContent,
	},
}

func generateSkill(baseDir string, t skillTarget, resources []resourceEntry) {
	dir := filepath.Join(baseDir, t.subdir)
	refDir := filepath.Join(dir, "references")
	if err := os.MkdirAll(refDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", refDir, err)
	}
	for name, content := range t.refs {
		if err := os.WriteFile(filepath.Join(refDir, name), []byte(content), 0o644); err != nil {
			log.Fatalf("write %s: %v", name, err)
		}
	}

	path := filepath.Join(dir, "SKILL.md")
	if err := writeSkill(path, resources, t); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}

	totalVerbs := 0
	for _, r := range resources {
		totalVerbs += len(r.verbs)
	}
	fmt.Printf("Generated skill: %d resources, %d verbs → %s\n", len(resources), totalVerbs, path)
}

func main() {
	skillDir := flag.String("skill-dir", ".", "base directory for generated skill")
	instructionDir := flag.String("instruction-dir", "./cmd/agent", "base directory for generated instruction and agent skill")
	flag.Parse()

	root := cmd.RootCmd
	root.InitDefaultHelpCmd()
	resources := collectResources(root)

	// External skill: includes setup.md for installation.
	generateSkill(*skillDir, extSkill, resources)

	// Agent instruction.
	if err := os.MkdirAll(*instructionDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", *instructionDir, err)
	}
	instrPath := filepath.Join(*instructionDir, "INSTRUCTION.md")
	if err := writeInstruction(instrPath); err != nil {
		log.Fatalf("write %s: %v", instrPath, err)
	}
	fmt.Printf("Generated instruction: → %s\n", instrPath)

	// Agent skill: embedded in the binary, no setup needed.
	generateSkill(*instructionDir, agentSkill, resources)
}
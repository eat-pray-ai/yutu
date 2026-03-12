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
	"github.com/spf13/pflag"

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

func main() {
	out := flag.String("out", "./skills", "output directory for generated skills")
	flag.Parse()

	root := cmd.RootCmd
	root.InitDefaultHelpCmd()

	for _, c := range root.Commands() {
		if !strings.HasPrefix(c.Short, "Manage") {
			continue
		}

		skill := c.Name()
		kebabSkill := camelToKebab(skill)
		dir := filepath.Join(*out, "youtube-"+kebabSkill)
		refDir := filepath.Join(dir, "references")

		if err := os.MkdirAll(refDir, 0o755); err != nil {
			log.Fatalf("mkdir %s: %v", refDir, err)
		}

		setupPath := filepath.Join(refDir, "setup.md")
		if err := writeSetup(setupPath); err != nil {
			log.Fatalf("write setup %s: %v", setupPath, err)
		}

		var verbs []*cobra.Command
		for _, sub := range c.Commands() {
			if sub.Name() == "help" {
				continue
			}
			verbs = append(verbs, sub)
		}
		sort.Slice(
			verbs, func(i, j int) bool {
				return verbs[i].Name() < verbs[j].Name()
			},
		)

		for _, verb := range verbs {
			refPath := filepath.Join(refDir, skill+"-"+verb.Name()+".md")
			if err := writeReference(refPath, c, verb); err != nil {
				log.Fatalf("write reference %s: %v", refPath, err)
			}
		}

		skillPath := filepath.Join(dir, "SKILL.md")
		if err := writeSkill(skillPath, c, verbs, kebabSkill); err != nil {
			log.Fatalf("write skill %s: %v", skillPath, err)
		}

		fmt.Printf("Generated skill: %s (%d verbs)\n", skill, len(verbs))
	}
}

func writeSetup(path string) error {
	return os.WriteFile(path, []byte(setupContent), 0o644)
}

// writeReference generates a reference Markdown file for a single verb
// command, including required-flag indicators and skill-oriented language.
func writeReference(path string, parent, verb *cobra.Command) error {
	var b strings.Builder

	skill := parent.Name()
	verbName := verb.Name()
	humanSkill := camelToWords(skill)
	title := titleCase(humanSkill) + " " + titleCase(verbName)

	b.WriteString(fmt.Sprintf("# %s\n\n", title))

	desc := rewriteToolPhrase(verb.Long)
	if desc != "" {
		b.WriteString(desc + "\n\n")
	}

	b.WriteString("## Usage\n\n")
	b.WriteString(
		fmt.Sprintf(
			"```bash\nyutu %s %s [flags]\n```\n", skill, verbName,
		),
	)

	requiredFlags := requiredFlagNames(verb)

	type flagEntry struct {
		name      string
		shorthand string
		usage     string
		required  bool
	}
	var flags []flagEntry
	verb.Flags().VisitAll(
		func(f *pflag.Flag) {
			if f.Name == "help" {
				return
			}
			usage := f.Usage
			if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
				usage += fmt.Sprintf(" (default %s)", formatDefault(f))
			}
			flags = append(
				flags, flagEntry{
					name:      f.Name,
					shorthand: f.Shorthand,
					usage:     escPipe(usage),
					required:  requiredFlags[f.Name],
				},
			)
		},
	)

	if len(flags) > 0 {
		b.WriteString("\n## Flags\n\n")
		b.WriteString("| Flag | Shorthand | Required | Description |\n")
		b.WriteString("|------|-----------|----------|-------------|\n")
		for _, f := range flags {
			sh := ""
			if f.shorthand != "" {
				sh = "`-" + f.shorthand + "`"
			}
			req := ""
			if f.required {
				req = "Yes"
			}
			b.WriteString(
				fmt.Sprintf(
					"| `--%-s` | %s | %s | %s |\n", f.name, sh, req, f.usage,
				),
			)
		}
	}

	if verb.Example != "" {
		b.WriteString("\n## Examples\n\n```bash\n")
		b.WriteString(strings.TrimSpace(verb.Example) + "\n")
		b.WriteString("```\n")
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// writeSkill generates the SKILL.md overview file for a resource command
// following skill-creator best practices: "pushy" description, overview
// paragraph, progressive-disclosure table, and quick-start snippet.
func writeSkill(
	path string, c *cobra.Command, verbs []*cobra.Command, kebabSkill string,
) error {
	var b strings.Builder

	skill := c.Name()
	humanSkill := camelToWords(skill)

	desc := buildDescription(c, verbs, humanSkill)
	skillName := "youtube-" + kebabSkill

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("name: %s\n", skillName))
	b.WriteString(
		fmt.Sprintf(
			"description: \"%s\"\n", strings.ReplaceAll(desc, `"`, `\"`),
		),
	)
	b.WriteString("compatibility: Requires the yutu CLI (brew install yutu), Google Cloud OAuth credentials (client_secret.json), and a cached OAuth token (youtube.token.json). Needs network access to the YouTube Data API.\n")
	b.WriteString("metadata:\n")
	b.WriteString("  author: eat-pray-ai\n")
	b.WriteString("---\n\n")

	b.WriteString(fmt.Sprintf("# YouTube %s\n\n", titleCase(humanSkill)))

	overview := rewriteToolPhrase(c.Long)
	if overview == "" {
		overview = c.Short
	}
	b.WriteString(overview + "\n\n")

	b.WriteString("## Before You Begin\n\n")
	b.WriteString("yutu requires Google Cloud Platform OAuth credentials and a cached token to access the YouTube API. ")
	b.WriteString("If you haven't set up yutu yet, read the [setup guide](references/setup.md) first.\n\n")

	b.WriteString("## Operations\n\n")
	b.WriteString("Read the linked reference for full flag details and examples.\n\n")

	if len(verbs) > 0 {
		b.WriteString("| Operation | Description | Reference |\n")
		b.WriteString("|-----------|-------------|----------|\n")
		for _, verb := range verbs {
			refFile := fmt.Sprintf("references/%s-%s.md", skill, verb.Name())
			b.WriteString(
				fmt.Sprintf(
					"| %s | %s | [details](%s) |\n",
					verb.Name(), escPipe(verb.Short), refFile,
				),
			)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Quick Start\n\n")
	b.WriteString(
		fmt.Sprintf(
			"```bash\n# Show all %s commands\nyutu %s --help\n", humanSkill, skill,
		),
	)
	for _, verb := range verbs {
		if verb.Name() == "list" {
			b.WriteString(
				fmt.Sprintf(
					"\n# List %s\nyutu %s list\n", humanSkill, skill,
				),
			)
			break
		}
	}
	b.WriteString("```\n")

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// buildDescription constructs a "pushy" skill description that encourages
// broad triggering, even when the user doesn't explicitly name the resource.
func buildDescription(
	c *cobra.Command, verbs []*cobra.Command, humanSkill string,
) string {
	base := rewriteToolPhrase(c.Long)
	if base == "" {
		base = c.Short
	}

	base += fmt.Sprintf(
		" Useful when working with YouTube %s — covers listing, creating, updating, and deleting %s via the yutu CLI. Includes setup and installation instructions for first-time users.",
		humanSkill, humanSkill,
	)

	base += " Triggers: " + naturalTriggerPhrases(humanSkill, verbs)

	return base
}

// naturalTriggerPhrases generates trigger phrases that feel like natural user
// requests rather than mechanical "verb + resource" pairs.
func naturalTriggerPhrases(humanSkill string, verbs []*cobra.Command) string {
	var phrases []string
	for _, v := range verbs {
		short := strings.ToLower(v.Short)
		name := v.Name()

		phrases = append(phrases, short)
		phrases = append(phrases, name+" "+humanSkill)
		phrases = append(phrases, name+" my "+humanSkill)
	}
	return strings.Join(phrases, ", ")
}

// requiredFlagNames returns the set of flag names marked as required on cmd.
func requiredFlagNames(cmd *cobra.Command) map[string]bool {
	required := make(map[string]bool)
	cmd.Flags().VisitAll(
		func(f *pflag.Flag) {
			if ann := f.Annotations; ann != nil {
				if _, ok := ann[cobra.BashCompOneRequiredFlag]; ok {
					required[f.Name] = true
				}
			}
		},
	)
	return required
}

// formatDefault returns a display-friendly default value for a flag.
func formatDefault(f *pflag.Flag) string {
	v := f.DefValue
	if f.Value.Type() == "string" && v != "" {
		return fmt.Sprintf("%q", v)
	}
	return v
}

// escPipe escapes pipe characters for markdown table cells.
func escPipe(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

// titleCase capitalises the first letter of each word in s.
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// rewriteToolPhrase replaces "Use this tool" phrasing with skill-oriented
// language so the text reads naturally inside a skill context.
func rewriteToolPhrase(s string) string {
	s = strings.Replace(s, "Use this tool to", "Use this skill to", 1)
	s = strings.Replace(s, "Use this tool", "Use this skill", 1)
	return strings.TrimSpace(s)
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

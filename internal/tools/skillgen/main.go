// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package main

import (
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
		dir := filepath.Join(*out, "yutu-"+skill)
		refDir := filepath.Join(dir, "references")

		if err := os.MkdirAll(refDir, 0o755); err != nil {
			log.Fatalf("mkdir %s: %v", refDir, err)
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
		if err := writeSkill(skillPath, c, verbs); err != nil {
			log.Fatalf("write skill %s: %v", skillPath, err)
		}

		fmt.Printf("Generated skill: %s (%d verbs)\n", skill, len(verbs))
	}
}

// writeReference generates a reference Markdown file for a single verb command.
func writeReference(path string, parent, verb *cobra.Command) error {
	var b strings.Builder

	skill := parent.Name()
	verbName := verb.Name()
	title := titleCase(skill) + " " + titleCase(verbName)

	b.WriteString(fmt.Sprintf("# %s Command\n\n", title))
	b.WriteString(stripToolPhrase(verb.Long) + "\n\n")
	b.WriteString("## Usage\n\n")
	b.WriteString(
		fmt.Sprintf(
			"```bash\nyutu %s %s [flags]\n```\n", skill, verbName,
		),
	)

	type flagEntry struct {
		name      string
		shorthand string
		usage     string
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
				},
			)
		},
	)

	if len(flags) > 0 {
		b.WriteString("\n## Flags\n\n")
		b.WriteString("| Flag | Shorthand | Description |\n")
		b.WriteString("|------|-----------|-------------|\n")
		for _, f := range flags {
			sh := ""
			if f.shorthand != "" {
				sh = "`-" + f.shorthand + "`"
			}
			b.WriteString(fmt.Sprintf("| `--%-s` | %s | %s |\n", f.name, sh, f.usage))
		}
	}

	if verb.Example != "" {
		b.WriteString("\n## Examples\n\n```bash\n")
		b.WriteString(strings.TrimSpace(verb.Example) + "\n")
		b.WriteString("```\n")
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// writeSkill generates the SKILL.md overview file for a resource command.
func writeSkill(path string, c *cobra.Command, verbs []*cobra.Command) error {
	var b strings.Builder

	skill := c.Name()
	desc := strings.Replace(
		c.Long, "Use this tool", "Use this skill when you need", 1,
	)
	if desc == "" {
		desc = c.Short
	}
	desc += " Triggers: " + triggerPhrases(skill, verbs)

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("name: yutu-%s\n", skill))
	b.WriteString(fmt.Sprintf("description: \"%s\"\n", strings.ReplaceAll(desc, `"`, `\"`)))
	b.WriteString("---\n\n")

	b.WriteString(fmt.Sprintf("# Yutu %s\n\n", titleCase(skill)))

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
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
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

// titleCase capitalises the first letter of a string.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func stripToolPhrase(s string) string {
	if i := strings.Index(s, " Use this tool"); i > 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

func triggerPhrases(skill string, verbs []*cobra.Command) string {
	var phrases []string
	human := camelToWords(skill)
	for _, v := range verbs {
		phrases = append(phrases, strings.ToLower(v.Short))
		phrases = append(phrases, v.Name()+" "+human)
	}
	return strings.Join(phrases, ", ")
}

func camelToWords(s string) string {
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
	return strings.Join(words, " ")
}

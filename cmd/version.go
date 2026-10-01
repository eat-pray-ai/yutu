// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"runtime/debug"
	"strings"

	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

const (
	verShort = "Show the version of yutu"
	verLong  = "Show the version of yutu."
	repo     = "Github/eat-pray-ai/yutu"
	repoUrl  = "https://github.com/eat-pray-ai/yutu"
)

var banner = strings.TrimPrefix(
	`
  _(\_/)_     yutu %s %s/%s
 │   ▶   │    The AI-powered toolkit for YouTube
 │ • x • │    📦 build: %s-%s
 ╰─/───\─╯    🌟 Star:  %s
`, "\n",
)

var (
	Version    = ""
	Commit     = ""
	CommitDate = ""
	Os         = ""
	Arch       = ""
	Builder    = "Gopher"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: verShort,
	Long:  verLong,
	Run: func(cmd *cobra.Command, _ []string) {
		info, ok := debug.ReadBuildInfo()
		if ok && Version == "" {
			Version = info.Main.Version

			settings := make(map[string]string)
			for _, setting := range info.Settings {
				settings[setting.Key] = setting.Value
			}

			if val, exists := settings["vcs.time"]; exists {
				CommitDate = val
			}
			if val, exists := settings["GOOS"]; exists {
				Os = val
			}
			if val, exists := settings["GOARCH"]; exists {
				Arch = val
			}
		}

		cmd.Printf(
			banner, Version, Os, Arch, Builder, CommitDate,
			termlink.Link(repo, repoUrl),
		)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

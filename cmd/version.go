// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"runtime/debug"

	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

const (
	verShort = "Show the version of yutu"
	verLong  = "Show the version of yutu"
	repo     = "Github/eat-pray-ai/yutu"
	repoUrl  = "https://github.com/eat-pray-ai/yutu"
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
	Run: func(cmd *cobra.Command, args []string) {
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

		cmd.Printf("🐰yutu %s %s/%s\n", Version, Os, Arch)
		cmd.Printf("📦build %s-%s\n", Builder, CommitDate)
		cmd.Printf("🌟Star: %s\n", termlink.Link(repo, repoUrl))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

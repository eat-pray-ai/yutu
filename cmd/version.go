// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

const (
	versionShort = "Show the version of yutu"
	versionLong  = "Show the version of yutu"
)

var (
	Version    = ""
	Commit     = ""
	CommitDate = ""
	Os         = ""
	Arch       = ""
	repo       = "https://github.com/eat-pray-ai/yutu"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: versionShort,
	Long:  versionLong,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("🐰yutu version %s %s/%s", Version, Os, Arch)
		if Commit != "" && CommitDate != "" {
			cmd.Printf("\n📦build %s-%s", Commit, CommitDate)
		}

		cmd.Println("\n🌟Star:", termlink.Link("Github/eat-pray-ai/yutu", repo))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

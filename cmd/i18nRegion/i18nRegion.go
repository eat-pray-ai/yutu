// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short      = "Manage YouTube i18n regions"
	long       = "Manage YouTube i18n regions. Use this tool to list available internationalization regions."
	hlUsage    = "Host language"
	regionURI  = "i18n://region/{hl}"
	regionName = "i18nRegions"
)

var (
	hl           string
	parts        []string
	output       string
	jsonpath     string
	defaultParts = []string{"id", "snippet"}
)

var i18nRegionCmd = &cobra.Command{
	Use:   "i18nRegion",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(i18nRegionCmd)
}

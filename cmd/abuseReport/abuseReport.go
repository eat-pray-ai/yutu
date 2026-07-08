// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package abuseReport

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/spf13/cobra"
)

const (
	short              = "Manage YouTube abuse reports"
	long               = "Manage YouTube abuse reports. Use this tool to report abusive content on YouTube."
	abuseTypesUsage    = "Abuse type IDs (e.g., spam, harassment)"
	descriptionUsage   = "Description of the abuse"
	subjectIdUsage     = "ID of the subject being reported"
	subjectTypeUsage   = "Type ID of the subject (e.g., video, comment, channel)"
	subjectUrlUsage    = "URL of the subject being reported"
	relatedEntityUsage = "ID of a related entity"
)

var (
	abuseTypes      []string
	description     string
	subjectId       string
	subjectTypeId   string
	subjectUrl      string
	relatedEntityId string
	parts           []string
)

var abuseReportCmd = &cobra.Command{
	Use:   "abuseReport",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(abuseReportCmd)
}

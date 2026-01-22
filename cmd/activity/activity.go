// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/spf13/cobra"
)

const (
	short     = "List YouTube activities"
	long      = "List YouTube activities, such as likes, favorites, uploads, etc"
	ciUsage   = "ID of the channel"
	homeUsage = "true|false|null"
	mineUsage = "true|false|null"
	paUsage   = "Filter on activities published after this date"
	pbUsage   = "Filter on activities published before this date"
	rcUsage   = ""
)

var (
	channelId       string
	home            = jsonschema.Ptr(false)
	maxResults      int64
	mine            = jsonschema.Ptr(false)
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
	output          string
	jsonpath        string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		boolMap := map[string]**bool{
			"home": &home,
			"mine": &mine,
		}
		utils.ResetBool(boolMap, cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(activityCmd)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	short     = "List YouTube activities"
	long      = "List YouTube activities, such as likes, favorites, uploads, etc\n\nExamples:\n  yutu activity list --mine\n  yutu activity list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10\n  yutu activity list --publishedAfter 2024-01-01T00:00:00Z --output json"
	ciUsage   = "ID of the channel"
	homeUsage = "true|false|null"
	mineUsage = "true|false|null"
	paUsage   = "Filter on activities published after this date"
	pbUsage   = "Filter on activities published before this date"
	rcUsage   = "Display the content as seen by viewers in this country"
)

var (
	channelId       string
	home            = new(false)
	maxResults      int64
	mine            = new(false)
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

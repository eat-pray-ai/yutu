// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	short       = "Manage YouTube live broadcasts"
	long        = "Manage YouTube live broadcasts. Use this tool to list, insert, update, delete, bind, transition, or insert cuepoints for live broadcasts."
	titleUsage  = "Title of the live broadcast"
	descUsage   = "Description of the live broadcast"
	psUsage     = "Privacy status (public, unlisted, private)"
	sstUsage    = "Scheduled start time in RFC 3339 format"
	setUsage    = "Scheduled end time in RFC 3339 format"
	bsUsage     = "Broadcast status filter (all, active, upcoming, completed)"
	btUsage     = "Broadcast type filter (all, event, persistent)"
	sidUsage    = "ID of the live stream to bind"
	ctUsage     = "Cue type (cueTypeAd)"
	cdsUsage    = "Cuepoint duration in seconds"
	ciomUsage   = "Cuepoint insertion offset time in milliseconds"
	cwmUsage    = "Cuepoint wall clock time in milliseconds"
	obococUsage = "Channel ID for content owner operations"
)

var (
	ids                           []string
	title                         string
	description                   string
	mine                          = new(false)
	broadcastStatus               string
	broadcastType                 string
	privacyStatus                 string
	scheduledStartTime            string
	scheduledEndTime              string
	streamId                      string
	cueType                       string
	cueDurationSecs               int64
	cueInsertionOffsetMs          int64
	cueWalltimeMs                 uint64
	maxResults                    int64
	parts                         []string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
)

var liveBroadcastCmd = &cobra.Command{
	Use:   "liveBroadcast",
	Short: short,
	Long:  long,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resetFlags(cmd.Flags())
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(liveBroadcastCmd)
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]**bool{
		"mine": &mine,
	}

	utils.ResetBool(boolMap, flagSet)
}

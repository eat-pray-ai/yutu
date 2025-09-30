// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

const (
	short        = "Manipulate YouTube channels"
	long         = "List or update YouTube channels"
	cidUsage     = "Return the channels within the specified guide category id"
	fhUsage      = "Return the channel associated with a YouTube handle"
	fuUsage      = "Return the channel associated with a YouTube username"
	hlUsage      = "Specifies the localization language of the metadata"
	mbmUsage     = "Return the channels managed by the authenticated user"
	mineUsage    = "Return the ids of channels owned by the authenticated user"
	msUsage      = "Return the channels subscribed to the authenticated user"
	countryUsage = "Country of the channel"
	curlUsage    = "Custom URL of the channel"
	dlUsage      = "The language of the channel's default title and description"
	descUsage    = "Description of the channel"
	titleUsage   = "Title of the channel"
)

var (
	categoryId      string
	forHandle       string
	forUsername     string
	hl              string
	ids             []string
	managedByMe     = utils.BoolPtr("false")
	maxResults      int64
	mine            = utils.BoolPtr("false")
	mySubscribers   = utils.BoolPtr("false")
	country         string
	customUrl       string
	defaultLanguage string
	description     string
	title           string
	parts           []string
	output          string
	jpath           string

	onBehalfOfContentOwner string
)

var channelCmd = &cobra.Command{
	Use:   "channel",
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
	cmd.RootCmd.AddCommand(channelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// channelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// channelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func resetFlags(flagSet *pflag.FlagSet) {
	boolMap := map[string]**bool{
		"managedByMe":   &managedByMe,
		"mine":          &mine,
		"mySubscribers": &mySubscribers,
	}

	utils.ResetBool(boolMap, flagSet)
}

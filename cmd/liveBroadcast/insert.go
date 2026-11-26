// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveBroadcast"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	insertTool    = "liveBroadcast-insert"
	insertConfirm = "Create live broadcast: %s"
	insertShort   = "Create a live broadcast"
	insertLong    = "Create a live broadcast. Use this tool to create a live broadcast for the authenticated user."
	insertExample = `  # Create a public live broadcast
  yutu liveBroadcast insert --title "My Broadcast" --privacyStatus public --scheduledStartTime 2026-01-01T00:00:00Z
  # Create a private broadcast with description
  yutu liveBroadcast insert --title "Test Stream" --description "Testing" --privacyStatus private --scheduledStartTime 2026-01-01T00:00:00Z`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"title", "scheduled_start_time"},
	Properties: map[string]*jsonschema.Schema{
		"title":                {Type: "string", Description: titleUsage},
		"description":          {Type: "string", Description: descUsage},
		"scheduled_start_time": {Type: "string", Description: sstUsage},
		"scheduled_end_time":   {Type: "string", Description: setUsage},
		"privacy_status": {
			Type: "string", Description: psUsage,
			Enum: []any{"public", "unlisted", "private"},
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["snippet","status","contentDetails"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input liveBroadcast.LiveBroadcast) string {
					return fmt.Sprintf(insertConfirm, input.Title)
				},
				func(input liveBroadcast.LiveBroadcast, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	title              string
	description        string
	scheduledStartTime string
	scheduledEndTime   string
	privacyStatus      string
	oboco              string
	obococ             string
	parts              []string
	output             string
	yes                bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.title, "title", "t", "", titleUsage)
	c.Flags().StringVarP(&o.description, "description", "d", "", descUsage)
	c.Flags().StringVarP(
		&o.scheduledStartTime, "scheduledStartTime", "S", "", sstUsage,
	)
	c.Flags().StringVarP(
		&o.scheduledEndTime, "scheduledEndTime", "E", "", setUsage,
	)
	c.Flags().StringVarP(&o.privacyStatus, "privacyStatus", "P", "", psUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet", "status", "contentDetails"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("title")
	_ = c.MarkFlagRequired("scheduledStartTime")
}

func (o *insertOpts) build(f *pflag.FlagSet) *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, liveBroadcast.WithDescription,
		),
		liveBroadcast.WithScheduledStartTime(o.scheduledStartTime),
		liveBroadcast.WithScheduledEndTime(o.scheduledEndTime),
		liveBroadcast.WithPrivacyStatus(o.privacyStatus),
		liveBroadcast.WithOBOCO(o.oboco),
		liveBroadcast.WithOBOCOC(o.obococ),
		liveBroadcast.WithParts(o.parts),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(c, opts.yes, fmt.Sprintf(insertConfirm, opts.title))
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

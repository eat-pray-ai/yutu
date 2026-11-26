// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json/jsontext"
	"fmt"
	"io"
	"strings"

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
	updateTool    = "liveBroadcast-update"
	updateConfirm = "Update live broadcast: %s"
	updateIdUsage = "ID of the live broadcast to update"
	updateShort   = "Update a live broadcast"
	updateLong    = "Update a live broadcast. Use this tool to update an existing live broadcast's settings."
	updateExample = `  # Update live broadcast title
  yutu liveBroadcast update --id broadcast123 --title "New Title"
  # Update privacy status
  yutu liveBroadcast update --id broadcast123 --privacyStatus unlisted`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
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
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input liveBroadcast.LiveBroadcast) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input liveBroadcast.LiveBroadcast, w io.Writer) error {
					liveBroadcast.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids                           []string
	title                         string
	description                   string
	scheduledStartTime            string
	scheduledEndTime              string
	privacyStatus                 string
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	output                        string
	yes                           bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
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
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithIds(o.ids),
		liveBroadcast.WithTitle(o.title),
		utils.OptIfChanged(
			f, "description", &o.description, liveBroadcast.WithDescription,
		),
		liveBroadcast.WithScheduledStartTime(o.scheduledStartTime),
		liveBroadcast.WithScheduledEndTime(o.scheduledEndTime),
		liveBroadcast.WithPrivacyStatus(o.privacyStatus),
		liveBroadcast.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		liveBroadcast.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		liveBroadcast.WithMaxResults(1),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

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
)

const (
	updateTool    = "liveBroadcast-update"
	updateConfirm = "Update live broadcast: %s"
	updateIdUsage = "ID of the live broadcast to update"
	updateShort   = "Update a live broadcast"
	updateLong    = "Update a live broadcast. Use this tool to update an existing live broadcast's settings."
	updateExample = `# Update live broadcast title
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
			Description: pkg.SilentUsage, Default: jsontext.Value(`"yaml"`),
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
					input.MaxResults = 1
					return input.Update(w)
				},
			),
		),
	)
	liveBroadcastCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringVarP(
		&scheduledStartTime, "scheduledStartTime", "S", "", sstUsage,
	)
	updateCmd.Flags().StringVarP(
		&scheduledEndTime, "scheduledEndTime", "E", "", setUsage,
	)
	updateCmd.Flags().StringVarP(&privacyStatus, "privacyStatus", "P", "", psUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	updateCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		return utils.ConfirmPreRun(
			c, fmt.Sprintf(updateConfirm, strings.Join(ids, ", ")),
		)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithIds(ids),
			liveBroadcast.WithTitle(title),
			liveBroadcast.WithDescription(description),
			liveBroadcast.WithScheduledStartTime(scheduledStartTime),
			liveBroadcast.WithScheduledEndTime(scheduledEndTime),
			liveBroadcast.WithPrivacyStatus(privacyStatus),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithMaxResults(1),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}

// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json"
	"io"

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
	insertTool    = "liveBroadcast-insert"
	insertShort   = "Insert a live broadcast"
	insertLong    = "Insert a live broadcast. Use this tool to create a new live broadcast for the authenticated user."
	insertExample = `# Create a public live broadcast
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
			Default: json.RawMessage(`["snippet","status","contentDetails"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
		}, cobramcp.GenToolHandler(
			insertTool,
			func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	insertCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	insertCmd.Flags().StringVarP(
		&scheduledStartTime, "scheduledStartTime", "S", "", sstUsage,
	)
	insertCmd.Flags().StringVarP(
		&scheduledEndTime, "scheduledEndTime", "E", "", setUsage,
	)
	insertCmd.Flags().StringVarP(&privacyStatus, "privacyStatus", "P", "", psUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		obococUsage,
	)
	insertCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet", "status", "contentDetails"},
		pkg.PartsUsage,
	)
	insertCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = insertCmd.MarkFlagRequired("title")
	_ = insertCmd.MarkFlagRequired("scheduledStartTime")
	cmd.AddMutationFlags(insertCmd)
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(c, "Would create live broadcast: %s", title)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveBroadcast.NewLiveBroadcast(
			liveBroadcast.WithTitle(title),
			liveBroadcast.WithDescription(description),
			liveBroadcast.WithScheduledStartTime(scheduledStartTime),
			liveBroadcast.WithScheduledEndTime(scheduledEndTime),
			liveBroadcast.WithPrivacyStatus(privacyStatus),
			liveBroadcast.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			liveBroadcast.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			liveBroadcast.WithParts(parts),
			liveBroadcast.WithOutput(output),
		)
		utils.HandleCmdError(input.Insert(c.OutOrStdout()), c)
	},
}

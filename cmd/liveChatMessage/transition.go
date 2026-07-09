// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"encoding/json"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatMessage"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	transitionTool    = "liveChatMessage-transition"
	transitionShort   = "Transition a live chat message"
	transitionLong    = "Transition a durable live chat event. Use this tool to change the status of a live chat message (e.g., close a poll)."
	transitionExample = `# Transition a live chat message to closed
yutu liveChatMessage transition --ids abc123 --status closed`
)

var transitionInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: "IDs of the live chat messages to transition",
			Items: &jsonschema.Schema{Type: "string"},
		},
		"status": {
			Type: "string", Enum: []any{"closed"},
			Description: statusUsage,
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
			Name: transitionTool, Title: transitionShort, Description: transitionLong,
			InputSchema: transitionInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			transitionTool,
			func(input liveChatMessage.LiveChatMessage, writer io.Writer) error {
				return input.Transition(writer)
			},
		),
	)
	liveChatMessageCmd.AddCommand(transitionCmd)

	transitionCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, "IDs of the messages to transition",
	)
	transitionCmd.Flags().StringVarP(&status, "status", "s", "", statusUsage)
	transitionCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)

	_ = transitionCmd.MarkFlagRequired("ids")
	_ = transitionCmd.MarkFlagRequired("status")
	cmd.AddMutationFlags(transitionCmd)
}

var transitionCmd = &cobra.Command{
	Use:     "transition",
	Short:   transitionShort,
	Long:    transitionLong,
	Example: transitionExample,
	Run: func(c *cobra.Command, args []string) {
		output, _ := c.Flags().GetString("output")
		err := cmd.Confirm(
			c, "Would transition live chat message(s) %s to status %s",
			strings.Join(ids, ", "), status,
		)
		if err != nil {
			utils.HandleCmdError(err, c)
			return
		}
		input := liveChatMessage.NewLiveChatMessage(
			liveChatMessage.WithIds(ids),
			liveChatMessage.WithStatus(status),
			liveChatMessage.WithOutput(output),
		)
		utils.HandleCmdError(input.Transition(c.OutOrStdout()), c)
	},
}

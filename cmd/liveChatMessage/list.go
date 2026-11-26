// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatMessage"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "liveChatMessage-list"
	listShort   = "List live chat messages"
	listLong    = "List live chat messages. Use this tool to list messages in a live chat."
	listExample = `  # List messages for a live chat
  yutu liveChatMessage list --liveChatId abc123
  # List messages with limit
  yutu liveChatMessage list --liveChatId abc123 --maxResults 10`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id": {Type: "string", Description: lcidUsage},
		"hl":           {Type: "string", Description: "Localization language for system messages"},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["snippet","authorDetails"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: jsontext.Value(`"json"`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cobramcp.GenToolHandler(
			listTool,
			func(input liveChatMessage.LiveChatMessage, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveChatMessageCmd.AddCommand(newListCmd())
}

type listOpts struct {
	liveChatId string
	hl         string
	maxResults int64
	parts      []string
	output     string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.liveChatId, "liveChatId", "L", "", lcidUsage)
	c.Flags().StringVarP(
		&o.hl, "hl", "l", "", "Localization language for system messages",
	)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet", "authorDetails"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
	_ = c.MarkFlagRequired("liveChatId")
}

func (o *listOpts) build() *liveChatMessage.LiveChatMessage {
	return liveChatMessage.NewLiveChatMessage(
		liveChatMessage.WithLiveChatId(o.liveChatId),
		liveChatMessage.WithHl(o.hl),
		liveChatMessage.WithMaxResults(o.maxResults),
		liveChatMessage.WithParts(o.parts),
		liveChatMessage.WithOutput(o.output),
	).(*liveChatMessage.LiveChatMessage)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

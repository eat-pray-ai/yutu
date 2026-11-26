// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveChatModerator"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "liveChatModerator-list"
	listShort   = "List live chat moderators"
	listLong    = "List live chat moderators. Use this tool to list moderators for a live chat."
	listExample = `  # List moderators for a live chat
  yutu liveChatModerator list --liveChatId abc123
  # List moderators with limit
  yutu liveChatModerator list --liveChatId abc123 --maxResults 10`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"live_chat_id"},
	Properties: map[string]*jsonschema.Schema{
		"live_chat_id": {Type: "string", Description: lcidUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["snippet"]`),
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
			func(input liveChatModerator.LiveChatModerator, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveChatModeratorCmd.AddCommand(newListCmd())
}

type listOpts struct {
	liveChatId string
	maxResults int64
	parts      []string
	output     string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.liveChatId, "liveChatId", "l", "", lcidUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
	_ = c.MarkFlagRequired("liveChatId")
}

func (o *listOpts) build() *liveChatModerator.LiveChatModerator {
	return liveChatModerator.NewLiveChatModerator(
		liveChatModerator.WithLiveChatId(o.liveChatId),
		liveChatModerator.WithMaxResults(o.maxResults),
		liveChatModerator.WithParts(o.parts),
		liveChatModerator.WithOutput(o.output),
	).(*liveChatModerator.LiveChatModerator)
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

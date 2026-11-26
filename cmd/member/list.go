// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/member"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "member-list"
	listShort   = "List channel members"
	listLong    = "List channel members. Use this tool to list information about channel members."
	listExample = `  # List all current members
  yutu member list --mode all_current
  # List members with limit
  yutu member list --maxResults 10`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"member_channel_id":   {Type: "string", Description: mcidUsage},
		"has_access_to_level": {Type: "string", Description: hatlUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"mode": {
			Type: "string", Description: mmUsage,
			Enum:    []any{"listMembersModeUnknown", "updates", "all_current"},
			Default: jsontext.Value(`"all_current"`),
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
			listTool, func(input member.Member, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	memberCmd.AddCommand(newListCmd())
}

type listOpts struct {
	memberChannelId  string
	hasAccessToLevel string
	maxResults       int64
	mode             string
	parts            []string
	output           string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&o.memberChannelId, "memberChannelId", "c", "", mcidUsage,
	)
	c.Flags().StringVarP(
		&o.hasAccessToLevel, "hasAccessToLevel", "a", "", hatlUsage,
	)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(&o.mode, "mode", "m", "all_current", mmUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *member.Member {
	return member.NewMember(
		member.WithMemberChannelId(o.memberChannelId),
		member.WithHasAccessToLevel(o.hasAccessToLevel),
		member.WithMaxResults(o.maxResults),
		member.WithMode(o.mode),
		member.WithParts(o.parts),
		member.WithOutput(o.output),
	).(*member.Member)
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

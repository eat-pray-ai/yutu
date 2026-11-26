// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

import (
	"encoding/json/jsontext"
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
	listTool     = "liveBroadcast-list"
	listIdsUsage = "Return live broadcasts with the given ids"
	listShort    = "List live broadcasts"
	listLong     = "List live broadcasts. Use this tool to list live broadcasts for the authenticated user."
	listExample  = `  # List my live broadcasts
  yutu liveBroadcast list --mine
  # List live broadcasts by ID
  yutu liveBroadcast list --ids broadcast1,broadcast2
  # List active live broadcasts
  yutu liveBroadcast list --mine --broadcastStatus active --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"mine":             {Type: "boolean", Description: "List broadcasts owned by the authenticated user"},
		"broadcast_status": {Type: "string", Enum: []any{"all", "active", "upcoming", "completed"}, Description: bsUsage},
		"broadcast_type":   {Type: "string", Enum: []any{"all", "event", "persistent"}, Description: btUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"), Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: pkg.OBOCOCUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","status"]`),
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
			listTool, func(input liveBroadcast.LiveBroadcast, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveBroadcastCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids             []string
	mine            bool
	broadcastStatus string
	broadcastType   string
	maxResults      int64
	oboco           string
	obococ          string
	parts           []string
	output          string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().BoolVarP(
		&o.mine, "mine", "m", false, "List broadcasts owned by the authenticated user",
	)
	c.Flags().StringVarP(
		&o.broadcastStatus, "broadcastStatus", "s", "", bsUsage,
	)
	c.Flags().StringVarP(&o.broadcastType, "broadcastType", "T", "", btUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B",
		"", pkg.OBOCOCUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build(f *pflag.FlagSet) *liveBroadcast.LiveBroadcast {
	return liveBroadcast.NewLiveBroadcast(
		liveBroadcast.WithIds(o.ids),
		utils.OptIfChanged(f, "mine", &o.mine, liveBroadcast.WithMine),
		liveBroadcast.WithBroadcastStatus(o.broadcastStatus),
		liveBroadcast.WithBroadcastType(o.broadcastType),
		liveBroadcast.WithMaxResults(o.maxResults),
		liveBroadcast.WithOBOCO(o.oboco),
		liveBroadcast.WithOBOCOC(o.obococ),
		liveBroadcast.WithParts(o.parts),
		liveBroadcast.WithOutput(o.output),
	).(*liveBroadcast.LiveBroadcast)
}

func newListCmd() *cobra.Command {
	var opts listOpts

	c := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Long:    listLong,
		Example: listExample,
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).List(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

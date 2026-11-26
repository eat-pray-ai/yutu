// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/liveStream"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	listTool     = "liveStream-list"
	listIdsUsage = "Return live streams with the given ids"
	listShort    = "List live streams"
	listLong     = "List live streams. Use this tool to list live streams for the authenticated user."
	listExample  = `  # List my live streams
  yutu liveStream list --mine
  # List live streams by ID
  yutu liveStream list --ids stream1,stream2
  # List live streams in JSON format
  yutu liveStream list --mine --output json --maxResults 10`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"mine": {Type: "boolean", Description: "List streams owned by the authenticated user"},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"), Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"on_behalf_of_content_owner_channel": {
			Type: "string", Description: obococUsage,
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","cdn","status"]`),
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
			listTool, func(input liveStream.LiveStream, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	liveStreamCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                           []string
	mine                          bool
	maxResults                    int64
	onBehalfOfContentOwner        string
	onBehalfOfContentOwnerChannel string
	parts                         []string
	output                        string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().BoolVarP(
		&o.mine, "mine", "m", false, "List streams owned by the authenticated user",
	)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B",
		"", obococUsage,
	)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "cdn", "status"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build(f *pflag.FlagSet) *liveStream.LiveStream {
	return liveStream.NewLiveStream(
		liveStream.WithIds(o.ids),
		utils.OptIfChanged(f, "mine", &o.mine, liveStream.WithMine),
		liveStream.WithMaxResults(o.maxResults),
		liveStream.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		liveStream.WithOnBehalfOfContentOwnerChannel(o.onBehalfOfContentOwnerChannel),
		liveStream.WithParts(o.parts),
		liveStream.WithOutput(o.output),
	).(*liveStream.LiveStream)
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

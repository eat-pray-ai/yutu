// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/videoAbuseReportReason"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "videoAbuseReportReason-list"
	listShort   = "List video abuse report reasons"
	listLong    = "List video abuse report reasons. Use this tool to list available abuse report reasons."
	listExample = `  # List video abuse report reasons
  yutu videoAbuseReportReason list --hl en`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"hl": {Type: "string", Description: hlUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet"]`),
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
			listTool, func(
				input videoAbuseReportReason.VideoAbuseReportReason, writer io.Writer,
			) error {
				return input.List(writer)
			},
		),
	)
	videoAbuseReportReasonCmd.AddCommand(newListCmd())
}

type listOpts struct {
	hl     string
	parts  []string
	output string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *videoAbuseReportReason.VideoAbuseReportReason {
	return videoAbuseReportReason.NewVideoAbuseReportReason(
		videoAbuseReportReason.WithHL(o.hl),
		videoAbuseReportReason.WithParts(o.parts),
		videoAbuseReportReason.WithOutput(o.output),
	).(*videoAbuseReportReason.VideoAbuseReportReason)
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/activity"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "activity-list"
	listShort   = "List activities"
	listLong    = "List activities. Use this tool to list activities such as uploads, likes, and favorites."
	listExample = `  # List my activities
  yutu activity list --for mine
  # List activities by channel ID with limit
  yutu activity list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --maxResults 10
  # List activities after a date in JSON format
  yutu activity list --publishedAfter 2024-01-01T00:00:00Z --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: ciUsage},
		"for": {
			Type: "string", Description: forUsage,
			Enum: []any{"home", "mine"},
		},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"), Minimum: new(float64(0)),
		},
		"published_after":  {Type: "string", Description: paUsage},
		"published_before": {Type: "string", Description: pbUsage},
		"region_code":      {Type: "string", Description: rcUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","contentDetails"]`),
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
			listTool, func(input activity.Activity, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	activityCmd.AddCommand(newListCmd())
}

type listOpts struct {
	channelId       string
	activityFor     string
	maxResults      int64
	publishedAfter  string
	publishedBefore string
	regionCode      string
	parts           []string
	output          string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", ciUsage)
	c.Flags().StringVarP(&o.activityFor, "for", "F", "", forUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.publishedAfter, "publishedAfter", "a", "", paUsage,
	)
	c.Flags().StringVarP(
		&o.publishedBefore, "publishedBefore", "b", "", pbUsage,
	)
	c.Flags().StringVarP(&o.regionCode, "regionCode", "r", "", rcUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "contentDetails"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *activity.Activity {
	return activity.NewActivity(
		activity.WithChannelId(o.channelId),
		activity.WithFor(o.activityFor),
		activity.WithMaxResults(o.maxResults),
		activity.WithPublishedAfter(o.publishedAfter),
		activity.WithPublishedBefore(o.publishedBefore),
		activity.WithRegionCode(o.regionCode),
		activity.WithParts(o.parts),
		activity.WithOutput(o.output),
	).(*activity.Activity)
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

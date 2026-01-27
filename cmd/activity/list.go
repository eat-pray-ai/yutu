// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/activity"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool = "activity-list"
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: ciUsage},
		"home":       {Type: "boolean", Description: homeUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"), Minimum: jsonschema.Ptr(float64(0)),
		},
		"mine":             {Type: "boolean", Description: mineUsage},
		"published_after":  {Type: "string", Description: paUsage},
		"published_before": {Type: "string", Description: pbUsage},
		"region_code":      {Type: "string", Description: rcUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet","contentDetails"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: short, Description: long,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, cmd.GenToolHandler(
			listTool, func(input activity.Activity, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	activityCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", ciUsage)
	listCmd.Flags().BoolVarP(home, "home", "H", true, homeUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&publishedAfter, "publishedAfter", "a", "", paUsage,
	)
	listCmd.Flags().StringVarP(
		&publishedBefore, "publishedBefore", "b", "", pbUsage,
	)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "contentDetails"},
		pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		input := activity.NewActivity(
			activity.WithChannelId(channelId),
			activity.WithHome(home),
			activity.WithMaxResults(maxResults),
			activity.WithMine(mine),
			activity.WithPublishedAfter(publishedAfter),
			activity.WithPublishedBefore(publishedBefore),
			activity.WithRegionCode(regionCode),
			activity.WithParts(parts),
			activity.WithOutput(output),
			activity.WithJsonpath(jsonpath),
		)
		if err := input.List(cmd.OutOrStdout()); err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

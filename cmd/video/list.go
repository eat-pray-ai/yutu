// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "video-list"
	listShort    = "List video's info"
	listLong     = "List video's info, such as title, description, etc\n\nExamples:\n  yutu video list --ids dQw4w9WgXcQ\n  yutu video list --ids dQw4w9WgXcQ,abc123 --output json\n  yutu video list --chart mostPopular --regionCode US --maxResults 10\n  yutu video list --myRating like --output yaml"
	listIdsUsage = "Return videos with the given ids"
	listMrUsage  = "Return videos liked/disliked by the authenticated user"
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"chart": {
			Type: "string", Description: chartUsage,
			Enum: []any{"chartUnspecified", "mostPopular", ""},
		},
		"hl":          {Type: "string", Description: hlUsage},
		"locale":      {Type: "string", Description: localUsage},
		"category_id": {Type: "string", Description: caidUsage},
		"region_code": {Type: "string", Description: rcUsage},
		"max_height": {
			Type: "number", Description: mhUsage, Minimum: new(float64(0)),
		},
		"max_width": {
			Type: "number", Description: mwUsage, Minimum: new(float64(0)),
		},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"), Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string"},
		"rating":                     {Type: "string", Description: listMrUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet","status"]`),
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
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cmd.GenToolHandler(
			listTool, func(input video.Video, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	videoCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&chart, "chart", "c", "", chartUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().StringVarP(&locale, "locale", "L", "", localUsage)
	listCmd.Flags().StringVarP(&categoryId, "videoCategoryId", "g", "", caidUsage)
	listCmd.Flags().StringVarP(&regionCode, "regionCode", "r", "", rcUsage)
	listCmd.Flags().Int64VarP(&maxHeight, "maxHeight", "H", 0, mhUsage)
	listCmd.Flags().Int64VarP(&maxWidth, "maxWidth", "W", 0, mwUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(&rating, "myRating", "R", "", listMrUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := video.NewVideo(
			video.WithIds(ids),
			video.WithChart(chart),
			video.WithHl(hl),
			video.WithLocale(locale),
			video.WithCategory(categoryId),
			video.WithRegionCode(regionCode),
			video.WithMaxHeight(maxHeight),
			video.WithMaxWidth(maxWidth),
			video.WithMaxResults(maxResults),
			video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			video.WithRating(rating),
			video.WithParts(parts),
			video.WithOutput(output),
			video.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

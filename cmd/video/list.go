// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "video-list"
	listIdsUsage = "Return videos with the given ids"
	listMrUsage  = "Return videos liked/disliked by the authenticated user"
	listShort    = "List videos"
	listLong     = "List videos. Use this tool to list videos."
	listExample  = `  # List a video by ID
  yutu video list --ids IBju0NwjQRc
  # List multiple videos in JSON format
  yutu video list --ids IBju0NwjQRc,wZ9z3FUuNvQ --output json
  # List most popular videos by region
  yutu video list --chart mostPopular --regionCode US --maxResults 10
  # List my liked videos
  yutu video list --myRating like --output yaml`
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
			Enum: []any{"chartUnspecified", "mostPopular"},
		},
		"hl":          {Type: "string", Description: hlUsage},
		"locale":      {Type: "string", Description: localeUsage},
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
			Default: jsontext.Value("5"), Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"rating":                     {Type: "string", Description: listMrUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet","status","statistics"]`),
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
			listTool, func(input video.Video, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	videoCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                    []string
	chart                  string
	hl                     string
	locale                 string
	categoryId             string
	regionCode             string
	maxHeight              int64
	maxWidth               int64
	maxResults             int64
	onBehalfOfContentOwner string
	rating                 string
	parts                  []string
	output                 string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, listIdsUsage)
	c.Flags().StringVarP(&o.chart, "chart", "c", "", chartUsage)
	c.Flags().StringVarP(&o.hl, "hl", "l", "", hlUsage)
	c.Flags().StringVarP(&o.locale, "locale", "L", "", localeUsage)
	c.Flags().StringVarP(&o.categoryId, "videoCategoryId", "g", "", caidUsage)
	c.Flags().StringVarP(&o.regionCode, "regionCode", "r", "", rcUsage)
	c.Flags().Int64VarP(&o.maxHeight, "maxHeight", "H", 0, mhUsage)
	c.Flags().Int64VarP(&o.maxWidth, "maxWidth", "W", 0, mwUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "",
		pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(&o.rating, "myRating", "R", "", listMrUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet", "status", "statistics"},
		pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *video.Video {
	return video.NewVideo(
		video.WithIds(o.ids),
		video.WithChart(o.chart),
		video.WithHl(o.hl),
		video.WithLocale(o.locale),
		video.WithCategory(o.categoryId),
		video.WithRegionCode(o.regionCode),
		video.WithMaxHeight(o.maxHeight),
		video.WithMaxWidth(o.maxWidth),
		video.WithMaxResults(o.maxResults),
		video.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
		video.WithRating(o.rating),
		video.WithParts(o.parts),
		video.WithOutput(o.output),
	).(*video.Video)
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

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "search-list"
	listShort   = "Search resources"
	listLong    = "Search resources. Use this tool to search for videos, channels, playlists, and other resources."
	listExample = `  # Search for videos by keyword
  yutu search list --q 'golang tutorial' --maxResults 10
  # Search videos by duration in JSON format
  yutu search list --q 'music' --types video --videoDuration medium --output json
  # List latest videos from a channel
  yutu search list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --order date --maxResults 5
  # Search for live streams
  yutu search list --q 'live coding' --eventType live --types video`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
		"channel_type": {
			Type: "string", Enum: []any{"channelTypeUnspecified", "any", "show"},
			Description: ctUsage, Default: jsontext.Value(`"channelTypeUnspecified"`),
		},
		"event_type": {
			Type: "string", Enum: []any{"none", "upcoming", "live", "completed"},
			Description: etUsage, Default: jsontext.Value(`"none"`),
		},
		"for": {
			Type: "string", Description: forUsage,
			Enum: []any{"mine", "developer", "contentOwner"},
		},
		"location":        {Type: "string", Description: locationUsage},
		"location_radius": {Type: "string", Description: lrUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"order": {
			Type: "string", Description: orderUsage,
			Default: jsontext.Value(`"relevance"`),
			Enum: []any{
				"searchSortUnspecified", "date", "rating",
				"viewCount", "relevance", "title", "videoCount",
			},
		},
		"published_after":    {Type: "string", Description: paUsage},
		"published_before":   {Type: "string", Description: pbUsage},
		"q":                  {Type: "string", Description: qUsage},
		"region_code":        {Type: "string", Description: rcUsage},
		"relevance_language": {Type: "string", Description: rlUsage},
		"safe_search": {
			Type: "string", Description: ssUsage,
			Enum: []any{
				"safeSearchSettingUnspecified", "none", "moderate", "strict",
			},
			Default: jsontext.Value(`"moderate"`),
		},
		"topic_id": {Type: "string", Description: tidUsage},
		"types": {
			Type: "array", Description: typesUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"video_caption": {
			Type: "string", Description: vcUsage,
			Enum: []any{
				"videoCaptionUnspecified", "any", "closedCaption", "none",
			},
			Default: jsontext.Value(`"any"`),
		},
		"video_category_id": {Type: "string", Description: vcidUsage},
		"video_definition": {
			Type: "string", Description: vdeUsage,
			Enum: []any{"any", "high", "standard"},
		},
		"video_dimension": {
			Type: "string", Enum: []any{"any", "2d", "3d"},
			Description: vdiUsage, Default: jsontext.Value(`"any"`),
		},
		"video_duration": {
			Type: "string", Description: vduUsage,
			Enum: []any{
				"videoDurationUnspecified", "any", "short", "medium", "long",
			},
			Default: jsontext.Value(`"any"`),
		},
		"video_embeddable": {
			Type: "string", Description: veUsage,
			Enum: []any{"videoEmbeddableUnspecified", "any", "true"},
		},
		"video_license": {
			Type: "string", Description: vlUsage,
			Enum: []any{"any", "youtube", "creativeCommon"},
		},
		"video_paid_product_placement": {
			Type: "string", Description: vpppUsage,
			Enum: []any{
				"videoPaidProductPlacementUnspecified", "any", "true",
			},
		},
		"video_syndicated": {
			Type: "string", Description: vsUsage,
			Enum: []any{"videoSyndicatedUnspecified", "any", "true"},
		},
		"video_type": {
			Type: "string", Description: vtUsage,
			Enum: []any{"videoTypeUnspecified", "any", "movie", "episode"},
		},
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
			listTool, func(input search.Search, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	searchCmd.AddCommand(newListCmd())
}

type listOpts struct {
	channelId                 string
	channelType               string
	eventType                 string
	searchFor                 string
	location                  string
	locationRadius            string
	maxResults                int64
	oboco                     string
	order                     string
	publishedAfter            string
	publishedBefore           string
	q                         string
	regionCode                string
	relevanceLanguage         string
	safeSearch                string
	topicId                   string
	types                     []string
	videoCaption              string
	videoCategoryId           string
	videoDefinition           string
	videoDimension            string
	videoDuration             string
	videoEmbeddable           string
	videoLicense              string
	videoPaidProductPlacement string
	videoSyndicated           string
	videoType                 string
	parts                     []string
	output                    string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVar(&o.channelId, "channelId", "", cidUsage)
	c.Flags().StringVar(
		&o.channelType, "channelType", "channelTypeUnspecified", ctUsage,
	)
	c.Flags().StringVar(&o.eventType, "eventType", "none", etUsage)
	c.Flags().StringVar(&o.searchFor, "for", "", forUsage)
	c.Flags().StringVar(&o.location, "location", "", locationUsage)
	c.Flags().StringVar(&o.locationRadius, "locationRadius", "", lrUsage)
	c.Flags().Int64Var(&o.maxResults, "maxResults", 5, pkg.MRUsage)
	c.Flags().StringVar(
		&o.oboco, "onBehalfOfContentOwner", "", pkg.OBOCOUsage,
	)
	c.Flags().StringVar(&o.order, "order", "relevance", orderUsage)
	c.Flags().StringVar(&o.publishedAfter, "publishedAfter", "", paUsage)
	c.Flags().StringVar(&o.publishedBefore, "publishedBefore", "", pbUsage)
	c.Flags().StringVar(&o.q, "q", "", qUsage)
	c.Flags().StringVar(&o.regionCode, "regionCode", "", rcUsage)
	c.Flags().StringVar(&o.relevanceLanguage, "relevanceLanguage", "", rlUsage)
	c.Flags().StringVar(&o.safeSearch, "safeSearch", "moderate", ssUsage)
	c.Flags().StringVar(&o.topicId, "topicId", "", tidUsage)
	c.Flags().StringSliceVar(&o.types, "types", nil, typesUsage)
	c.Flags().StringVar(&o.videoCaption, "videoCaption", "any", vcUsage)
	c.Flags().StringVar(&o.videoCategoryId, "videoCategoryId", "", vcidUsage)
	c.Flags().StringVar(&o.videoDefinition, "videoDefinition", "", vdeUsage)
	c.Flags().StringVar(&o.videoDimension, "videoDimension", "any", vdiUsage)
	c.Flags().StringVar(&o.videoDuration, "videoDuration", "any", vduUsage)
	c.Flags().StringVar(&o.videoEmbeddable, "videoEmbeddable", "", veUsage)
	c.Flags().StringVar(&o.videoLicense, "videoLicense", "", vlUsage)
	c.Flags().StringVar(
		&o.videoPaidProductPlacement, "videoPaidProductPlacement", "", vpppUsage,
	)
	c.Flags().StringVar(&o.videoSyndicated, "videoSyndicated", "", vsUsage)
	c.Flags().StringVar(&o.videoType, "videoType", "", vtUsage)
	c.Flags().StringSliceVar(
		&o.parts, "parts", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVar(&o.output, "output", "table", pkg.TableUsage)
}

func (o *listOpts) build() *search.Search {
	return search.NewSearch(
		search.WithChannelId(o.channelId),
		search.WithChannelType(o.channelType),
		search.WithEventType(o.eventType),
		search.WithFor(o.searchFor),
		search.WithLocation(o.location),
		search.WithLocationRadius(o.locationRadius),
		search.WithMaxResults(o.maxResults),
		search.WithOBOCO(o.oboco),
		search.WithOrder(o.order),
		search.WithPublishedAfter(o.publishedAfter),
		search.WithPublishedBefore(o.publishedBefore),
		search.WithQ(o.q),
		search.WithRegionCode(o.regionCode),
		search.WithRelevanceLanguage(o.relevanceLanguage),
		search.WithSafeSearch(o.safeSearch),
		search.WithTopicId(o.topicId),
		search.WithTypes(o.types),
		search.WithVideoCaption(o.videoCaption),
		search.WithVideoCategoryId(o.videoCategoryId),
		search.WithVideoDefinition(o.videoDefinition),
		search.WithVideoDimension(o.videoDimension),
		search.WithVideoDuration(o.videoDuration),
		search.WithVideoEmbeddable(o.videoEmbeddable),
		search.WithVideoLicense(o.videoLicense),
		search.WithVideoPaidProductPlacement(o.videoPaidProductPlacement),
		search.WithVideoSyndicated(o.videoSyndicated),
		search.WithVideoType(o.videoType),
		search.WithParts(o.parts),
		search.WithOutput(o.output),
	).(*search.Search)
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

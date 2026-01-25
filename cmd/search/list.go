// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool = "search-list"
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
		"channel_type": {
			Type: "string", Enum: []any{"channelTypeUnspecified", "any", "show"},
			Description: ctUsage, Default: json.RawMessage(`"channelTypeUnspecified"`),
		},
		"event_type": {
			Type: "string", Enum: []any{"none", "upcoming", "live", "completed"},
			Description: etUsage, Default: json.RawMessage(`"none"`),
		},
		"for_content_owner": {Type: "boolean", Description: fcoUsage},
		"for_developer":     {Type: "boolean", Description: fdUsage},
		"for_mine":          {Type: "boolean", Description: fmUsage},
		"location":          {Type: "string", Description: locationUsage},
		"location_radius":   {Type: "string", Description: lrUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"on_behalf_of_content_owner": {Type: "string"},
		"order": {
			Type: "string", Description: orderUsage,
			Default: json.RawMessage(`"relevance"`),
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
			Default: json.RawMessage(`"moderate"`),
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
			Default: json.RawMessage(`"any"`),
		},
		"video_category_id": {Type: "string", Description: vcidUsage},
		"video_definition":  {Type: "string", Description: vdeUsage},
		"video_dimension": {
			Type: "string", Enum: []any{"any", "2d", "3d"},
			Description: vdiUsage, Default: json.RawMessage(`"any"`),
		},
		"video_duration": {
			Type: "string", Description: vduUsage,
			Enum: []any{
				"videoDurationUnspecified", "any", "short", "medium", "long",
			},
			Default: json.RawMessage(`"any"`),
		},
		"video_embeddable": {
			Type: "string", Description: veUsage,
			Enum: []any{"videoEmbeddableUnspecified", "any", "true", ""},
		},
		"video_license": {
			Type: "string", Description: vlUsage,
			Enum: []any{"any", "youtube", "creativeCommon", ""},
		},
		"video_paid_product_placement": {
			Type: "string", Description: vpppUsage,
			Enum: []any{
				"videoPaidProductPlacementUnspecified", "any", "true", "",
			},
		},
		"video_syndicated": {
			Type: "string", Description: vsUsage,
			Enum: []any{"videoSyndicatedUnspecified", "any", "true", ""},
		},
		"video_type": {
			Type: "string", Description: vtUsage,
			Enum: []any{"videoTypeUnspecified", "any", "movie", "episode", ""},
		},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet"]`),
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
		}, listHandler,
	)
	searchCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&channelId, "channelId", "", cidUsage)
	listCmd.Flags().StringVar(
		&channelType, "channelType", "channelTypeUnspecified", ctUsage,
	)
	listCmd.Flags().StringVar(&eventType, "eventType", "none", etUsage)
	listCmd.Flags().BoolVar(
		forContentOwner, "forContentOwner", false, fcoUsage,
	)
	listCmd.Flags().BoolVar(forDeveloper, "forDeveloper", false, fdUsage)
	listCmd.Flags().BoolVar(forMine, "forMine", false, fmUsage)
	listCmd.Flags().StringVar(&location, "location", "", locationUsage)
	listCmd.Flags().StringVar(&locationRadius, "locationRadius", "", lrUsage)
	listCmd.Flags().Int64Var(&maxResults, "maxResults", 5, pkg.MRUsage)
	listCmd.Flags().StringVar(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "", "",
	)
	listCmd.Flags().StringVar(&order, "order", "relevance", orderUsage)
	listCmd.Flags().StringVar(&publishedAfter, "publishedAfter", "", paUsage)
	listCmd.Flags().StringVar(&publishedBefore, "publishedBefore", "", pbUsage)
	listCmd.Flags().StringVar(&q, "q", "", qUsage)
	listCmd.Flags().StringVar(&regionCode, "regionCode", "", rcUsage)
	listCmd.Flags().StringVar(&relevanceLanguage, "relevanceLanguage", "", rlUsage)
	listCmd.Flags().StringVar(&safeSearch, "safeSearch", "moderate", ssUsage)
	listCmd.Flags().StringVar(&topicId, "topicId", "", tidUsage)
	listCmd.Flags().StringSliceVar(&types, "types", []string{}, typesUsage)
	listCmd.Flags().StringVar(&videoCaption, "videoCaption", "any", vcUsage)
	listCmd.Flags().StringVar(&videoCategoryId, "videoCategoryId", "", vcidUsage)
	listCmd.Flags().StringVar(&videoDefinition, "videoDefinition", "", vdeUsage)
	listCmd.Flags().StringVar(&videoDimension, "videoDimension", "any", vdiUsage)
	listCmd.Flags().StringVar(&videoDuration, "videoDuration", "any", vduUsage)
	listCmd.Flags().StringVar(&videoEmbeddable, "videoEmbeddable", "", veUsage)
	listCmd.Flags().StringVar(&videoLicense, "videoLicense", "", vlUsage)
	listCmd.Flags().StringVar(
		&videoPaidProductPlacement, "videoPaidProductPlacement", "", vpppUsage,
	)
	listCmd.Flags().StringVar(&videoSyndicated, "videoSyndicated", "", vsUsage)
	listCmd.Flags().StringVar(&videoType, "videoType", "", vtUsage)
	listCmd.Flags().StringSliceVar(
		&parts, "parts", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVar(&output, "output", "table", pkg.TableUsage)
	listCmd.Flags().StringVar(&jsonpath, "jsonpath", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		input := search.NewSearch(
			search.WithChannelId(channelId),
			search.WithChannelType(channelType),
			search.WithEventType(eventType),
			search.WithForContentOwner(forContentOwner),
			search.WithForDeveloper(forDeveloper),
			search.WithForMine(forMine),
			search.WithLocation(location),
			search.WithLocationRadius(locationRadius),
			search.WithMaxResults(maxResults),
			search.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			search.WithOrder(order),
			search.WithPublishedAfter(publishedAfter),
			search.WithPublishedBefore(publishedBefore),
			search.WithQ(q),
			search.WithRegionCode(regionCode),
			search.WithRelevanceLanguage(relevanceLanguage),
			search.WithSafeSearch(safeSearch),
			search.WithTopicId(topicId),
			search.WithTypes(types),
			search.WithVideoCaption(videoCaption),
			search.WithVideoCategoryId(videoCategoryId),
			search.WithVideoDefinition(videoDefinition),
			search.WithVideoDimension(videoDimension),
			search.WithVideoDuration(videoDuration),
			search.WithVideoEmbeddable(videoEmbeddable),
			search.WithVideoLicense(videoLicense),
			search.WithVideoPaidProductPlacement(videoPaidProductPlacement),
			search.WithVideoSyndicated(videoSyndicated),
			search.WithVideoType(videoType),
			search.WithParts(parts),
			search.WithOutput(output),
			search.WithJsonpath(jsonpath),
		)

		err := input.List(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func listHandler(
	ctx context.Context, req *mcp.CallToolRequest, input search.Search,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: listTool, MinInterval: time.Second},
		),
	)

	var writer bytes.Buffer
	err := input.List(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

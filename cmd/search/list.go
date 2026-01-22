// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/search"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool = "search-list"
)

type listIn struct {
	ChannelId                 string   `json:"channelId"`
	ChannelType               string   `json:"channelType"`
	EventType                 string   `json:"eventType"`
	ForContentOwner           *string  `json:"forContentOwner,omitempty"`
	ForDeveloper              *string  `json:"forDeveloper,omitempty"`
	ForMine                   *string  `json:"forMine,omitempty"`
	Location                  string   `json:"location"`
	LocationRadius            string   `json:"locationRadius"`
	MaxResults                int64    `json:"maxResults"`
	OnBehalfOfContentOwner    string   `json:"onBehalfOfContentOwner"`
	Order                     string   `json:"order"`
	PublishedAfter            string   `json:"publishedAfter"`
	PublishedBefore           string   `json:"publishedBefore"`
	Q                         string   `json:"q"`
	RegionCode                string   `json:"regionCode"`
	RelevanceLanguage         string   `json:"relevanceLanguage"`
	SafeSearch                string   `json:"safeSearch"`
	TopicId                   string   `json:"topicId"`
	Types                     []string `json:"types"`
	VideoCaption              string   `json:"videoCaption"`
	VideoCategoryId           string   `json:"videoCategoryId"`
	VideoDefinition           string   `json:"videoDefinition"`
	VideoDimension            string   `json:"videoDimension"`
	VideoDuration             string   `json:"videoDuration"`
	VideoEmbeddable           string   `json:"videoEmbeddable"`
	VideoLicense              string   `json:"videoLicense"`
	VideoPaidProductPlacement string   `json:"videoPaidProductPlacement"`
	VideoSyndicated           string   `json:"videoSyndicated"`
	VideoType                 string   `json:"videoType"`
	Parts                     []string `json:"parts"`
	Output                    string   `json:"output"`
	Jsonpath                  string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channelId": {
			Type: "string", Description: cidUsage,
			Default: json.RawMessage(`""`),
		},
		"channelType": {
			Type: "string", Enum: []any{"channelTypeUnspecified", "any", "show"},
			Description: ctUsage, Default: json.RawMessage(`"channelTypeUnspecified"`),
		},
		"eventType": {
			Type: "string", Enum: []any{"none", "upcoming", "live", "completed"},
			Description: etUsage, Default: json.RawMessage(`"none"`),
		},
		"forContentOwner": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: fcoUsage, Default: json.RawMessage(`""`),
		},
		"forDeveloper": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: fdUsage, Default: json.RawMessage(`""`),
		},
		"forMine": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: fmUsage, Default: json.RawMessage(`""`),
		},
		"location": {
			Type: "string", Description: locationUsage,
			Default: json.RawMessage(`""`),
		},
		"locationRadius": {
			Type: "string", Description: lrUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"order": {
			Type: "string", Description: orderUsage,
			Default: json.RawMessage(`"relevance"`),
		},
		"publishedAfter": {
			Type: "string", Description: paUsage,
			Default: json.RawMessage(`""`),
		},
		"publishedBefore": {
			Type: "string", Description: pbUsage,
			Default: json.RawMessage(`""`),
		},
		"q": {
			Type: "string", Description: qUsage,
			Default: json.RawMessage(`""`),
		},
		"regionCode": {
			Type: "string", Description: rcUsage,
			Default: json.RawMessage(`""`),
		},
		"relevanceLanguage": {
			Type: "string", Description: rlUsage,
			Default: json.RawMessage(`""`),
		},
		"safeSearch": {
			Type: "string",
			Enum: []any{
				"safeSearchSettingUnspecified", "none", "moderate", "strict",
			},
			Description: ssUsage, Default: json.RawMessage(`"moderate"`),
		},
		"topicId": {
			Type: "string", Description: tidUsage,
			Default: json.RawMessage(`""`),
		},
		"types": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: typesUsage,
			Default:     json.RawMessage(`[]`),
		},
		"videoCaption": {
			Type: "string",
			Enum: []any{
				"videoCaptionUnspecified", "any", "closedCaption", "none",
			},
			Description: vcUsage, Default: json.RawMessage(`"any"`),
		},
		"videoCategoryId": {
			Type: "string", Description: vcidUsage,
			Default: json.RawMessage(`""`),
		},
		"videoDefinition": {
			Type: "string", Description: vdeUsage,
			Default: json.RawMessage(`""`),
		},
		"videoDimension": {
			Type: "string", Enum: []any{"any", "2d", "3d"},
			Description: vdiUsage, Default: json.RawMessage(`"any"`),
		},
		"videoDuration": {
			Type: "string",
			Enum: []any{
				"videoDurationUnspecified", "any", "short", "medium", "long",
			},
			Description: vduUsage, Default: json.RawMessage(`"any"`),
		},
		"videoEmbeddable": {
			Type:        "string",
			Enum:        []any{"videoEmbeddableUnspecified", "any", "true", ""},
			Description: veUsage, Default: json.RawMessage(`""`),
		},
		"videoLicense": {
			Type: "string", Enum: []any{"any", "youtube", "creativeCommon", ""},
			Description: vlUsage, Default: json.RawMessage(`""`),
		},
		"videoPaidProductPlacement": {
			Type: "string",
			Enum: []any{
				"videoPaidProductPlacementUnspecified", "any", "true", "",
			},
			Description: vpppUsage, Default: json.RawMessage(`""`),
		},
		"videoSyndicated": {
			Type:        "string",
			Enum:        []any{"videoSyndicatedUnspecified", "any", "true", ""},
			Description: vsUsage, Default: json.RawMessage(`""`),
		},
		"videoType": {
			Type:        "string",
			Enum:        []any{"videoTypeUnspecified", "any", "movie", "episode", ""},
			Description: vtUsage, Default: json.RawMessage(`""`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet"]`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "table"},
			Description: pkg.TableUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
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
		err := list(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func listHandler(
	ctx context.Context, req *mcp.CallToolRequest, input listIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{LoggerName: listTool, MinInterval: time.Second},
		),
	)

	channelId = input.ChannelId
	channelType = input.ChannelType
	eventType = input.EventType
	forContentOwner = utils.BoolPtr(*input.ForContentOwner)
	forDeveloper = utils.BoolPtr(*input.ForDeveloper)
	forMine = utils.BoolPtr(*input.ForMine)
	location = input.Location
	locationRadius = input.LocationRadius
	maxResults = input.MaxResults
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	order = input.Order
	publishedAfter = input.PublishedAfter
	publishedBefore = input.PublishedBefore
	q = input.Q
	regionCode = input.RegionCode
	relevanceLanguage = input.RelevanceLanguage
	safeSearch = input.SafeSearch
	topicId = input.TopicId
	types = input.Types
	videoCaption = input.VideoCaption
	videoCategoryId = input.VideoCategoryId
	videoDefinition = input.VideoDefinition
	videoDimension = input.VideoDimension
	videoDuration = input.VideoDuration
	videoEmbeddable = input.VideoEmbeddable
	videoLicense = input.VideoLicense
	videoPaidProductPlacement = input.VideoPaidProductPlacement
	videoSyndicated = input.VideoSyndicated
	videoType = input.VideoType
	parts = input.Parts
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	s := search.NewSearch(
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
		search.WithService(nil),
	)

	return s.List(parts, output, jsonpath, writer)
}

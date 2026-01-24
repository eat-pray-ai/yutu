// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/video"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "video-list"
	listShort    = "List video's info"
	listLong     = "List video's info, such as title, description, etc"
	listIdsUsage = "Return videos with the given ids"
	listMrUsage  = "Return videos liked/disliked by the authenticated user"
)

type listIn struct {
	Ids                    []string `json:"ids"`
	Chart                  string   `json:"chart"`
	Hl                     string   `json:"hl"`
	Locale                 string   `json:"locale"`
	VideoCategoryId        string   `json:"videoCategoryId"`
	RegionCode             string   `json:"regionCode"`
	MaxHeight              int64    `json:"maxHeight"`
	MaxWidth               int64    `json:"maxWidth"`
	MaxResults             int64    `json:"maxResults"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	MyRating               string   `json:"myRating"`
	Parts                  []string `json:"parts"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"chart": {
			Type: "string", Enum: []any{"chartUnspecified", "mostPopular", ""},
			Description: chartUsage, Default: json.RawMessage(`""`),
		},
		"hl": {
			Type: "string", Description: hlUsage,
			Default: json.RawMessage(`""`),
		},
		"locale": {
			Type: "string", Description: localUsage,
			Default: json.RawMessage(`""`),
		},
		"videoCategoryId": {
			Type: "string", Description: caidUsage,
			Default: json.RawMessage(`""`),
		},
		"regionCode": {
			Type: "string", Description: rcUsage,
			Default: json.RawMessage(`""`),
		},
		"maxHeight": {
			Type: "number", Description: mhUsage,
			Default: json.RawMessage("0"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"maxWidth": {
			Type: "number", Description: mwUsage,
			Default: json.RawMessage("0"),
			Minimum: jsonschema.Ptr(float64(0)),
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
		"myRating": {
			Type: "string", Description: listMrUsage,
			Default: json.RawMessage(`""`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["id","snippet","status"]`),
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
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, listHandler,
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

	ids = input.Ids
	chart = input.Chart
	hl = input.Hl
	locale = input.Locale
	categoryId = input.VideoCategoryId
	regionCode = input.RegionCode
	maxHeight = input.MaxHeight
	maxWidth = input.MaxWidth
	maxResults = input.MaxResults
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	rating = input.MyRating
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
	v := video.NewVideo(
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
		video.WithService(nil),
	)

	return v.List(parts, output, jsonpath, writer)
}

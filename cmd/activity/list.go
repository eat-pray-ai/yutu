// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

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

type listIn struct {
	ChannelId       string   `json:"channelId"`
	Home            *bool    `json:"home,omitempty"`
	MaxResults      int64    `json:"maxResults"`
	Mine            *bool    `json:"mine,omitempty"`
	PublishedAfter  string   `json:"publishedAfter"`
	PublishedBefore string   `json:"publishedBefore"`
	RegionCode      string   `json:"regionCode"`
	Parts           []string `json:"parts"`
	Output          string   `json:"output"`
	Jsonpath        string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"channelId": {Type: "string", Description: ciUsage},
		"home": {
			Type: "boolean", Enum: []any{true, false}, Description: homeUsage,
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"), Minimum: jsonschema.Ptr(float64(0)),
		},
		"mine": {
			Type: "boolean", Enum: []any{true, false}, Description: mineUsage,
		},
		"publishedAfter":  {Type: "string", Description: paUsage},
		"publishedBefore": {Type: "string", Description: pbUsage},
		"regionCode":      {Type: "string", Description: rcUsage},
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
		}, listHandler,
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
		input := &listIn{
			ChannelId:       channelId,
			Home:            home,
			MaxResults:      maxResults,
			Mine:            mine,
			PublishedAfter:  publishedAfter,
			PublishedBefore: publishedAfter,
			RegionCode:      regionCode,
			Parts:           parts,
			Output:          output,
			Jsonpath:        jsonpath,
		}
		if err := input.call(cmd.OutOrStdout()); err != nil {
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

	var writer bytes.Buffer
	if err := input.call(&writer); err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func (l *listIn) call(writer io.Writer, opts ...activity.Option) error {
	defaultOpts := []activity.Option{
		activity.WithChannelId(l.ChannelId),
		activity.WithHome(l.Home),
		activity.WithMaxResults(l.MaxResults),
		activity.WithMine(l.Mine),
		activity.WithPublishedAfter(l.PublishedAfter),
		activity.WithPublishedBefore(l.PublishedBefore),
		activity.WithRegionCode(l.RegionCode),
		activity.WithService(nil),
	}
	defaultOpts = append(defaultOpts, opts...)
	a := activity.NewActivity(defaultOpts...)

	return a.List(l.Parts, l.Output, l.Jsonpath, writer)
}

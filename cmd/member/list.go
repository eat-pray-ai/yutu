// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/member"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool = "member-list"
)

type listIn struct {
	MemberChannelId  string   `json:"memberChannelId"`
	HasAccessToLevel string   `json:"hasAccessToLevel"`
	MaxResults       int64    `json:"maxResults"`
	Mode             string   `json:"mode"`
	Parts            []string `json:"parts"`
	Output           string   `json:"output"`
	Jsonpath         string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"memberChannelId": {
			Type: "string", Description: mcidUsage,
			Default: json.RawMessage(`""`),
		},
		"hasAccessToLevel": {
			Type: "string", Description: hatlUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"mode": {
			Type:        "string",
			Enum:        []any{"listMembersModeUnknown", "updates", "all_current"},
			Description: mmUsage, Default: json.RawMessage(`"all_current"`),
		},
		"parts": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: pkg.PartsUsage,
			Default:     json.RawMessage(`["snippet"]`),
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
	memberCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&memberChannelId, "memberChannelId", "c", "", mcidUsage,
	)
	listCmd.Flags().StringVarP(
		&hasAccessToLevel, "hasAccessToLevel", "a", "", hatlUsage,
	)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(&mode, "mode", "m", "all_current", mmUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
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

	memberChannelId = input.MemberChannelId
	hasAccessToLevel = input.HasAccessToLevel
	maxResults = input.MaxResults
	mode = input.Mode
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
	m := member.NewMember(
		member.WithMemberChannelId(memberChannelId),
		member.WithHasAccessToLevel(hasAccessToLevel),
		member.WithMaxResults(maxResults),
		member.WithMode(mode),
		member.WithService(nil),
	)

	return m.List(parts, output, jsonpath, writer)
}

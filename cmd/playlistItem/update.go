// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlistItem"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "playlistItem-update"
	updateShort   = "Update a playlist item"
	updateLong    = "Update a playlist item's info, such as title, description, etc"
	updateIdUsage = "ID of the playlist item to update"
)

type updateIn struct {
	Ids                    []string `json:"ids"`
	Title                  string   `json:"title"`
	Description            string   `json:"description"`
	Privacy                string   `json:"privacy"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: updateIdUsage,
			Default:     json.RawMessage(`[]`),
		},
		"title": {
			Type: "string", Description: titleUsage,
			Default: json.RawMessage(`""`),
		},
		"description": {
			Type: "string", Description: descUsage,
			Default: json.RawMessage(`""`),
		},
		"privacy": {
			Type: "string", Enum: []any{"public", "private", "unlisted", ""},
			Description: privacyUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
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
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	playlistItemCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&description, "description", "d", "", descUsage)
	updateCmd.Flags().StringVarP(&privacy, "privacy", "p", "", privacyUsage)
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, req *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: updateTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.Ids
	title = input.Title
	description = input.Description
	privacy = input.Privacy
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func update(writer io.Writer) error {
	pi := playlistItem.NewPlaylistItem(
		playlistItem.WithIDs(ids),
		playlistItem.WithTitle(title),
		playlistItem.WithDescription(description),
		playlistItem.WithPrivacy(privacy),
		playlistItem.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlistItem.WithMaxResults(1),
		playlistItem.WithService(nil),
	)

	return pi.Update(output, jsonpath, writer)
}

// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "channel-update"
	updateShort   = "Update channel's info"
	updateLong    = "Update channel's info, such as title, description, etc"
	updateIdUsage = "ID of the channel to update"
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"country":          {Type: "string", Description: countryUsage},
		"custom_url":       {Type: "string", Description: curlUsage},
		"default_language": {Type: "string", Description: dlUsage},
		"description":      {Type: "string", Description: descUsage},
		"title":            {Type: "string", Description: titleUsage},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	channelCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().StringVarP(&country, "country", "c", "", countryUsage)
	updateCmd.Flags().StringVarP(&customUrl, "customUrl", "u", "", curlUsage)
	updateCmd.Flags().StringVarP(
		&defaultLanguage, "defaultLanguage", "l", "", dlUsage,
	)
	updateCmd.Flags().StringVarP(
		&description, "description", "d", "", descUsage,
	)
	updateCmd.Flags().StringVarP(&title, "title", "t", "", titleUsage)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := channel.NewChannel(
			channel.WithIDs(ids),
			channel.WithCountry(country),
			channel.WithCustomUrl(customUrl),
			channel.WithDefaultLanguage(defaultLanguage),
			channel.WithDescription(description),
			channel.WithTitle(title),
			channel.WithOutput(output),
			channel.WithJsonpath(jsonpath),
		)
		err := input.Update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, req *mcp.CallToolRequest, input channel.Channel,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: updateTool, MinInterval: time.Second,
			},
		),
	)

	var writer bytes.Buffer
	err := input.Update(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

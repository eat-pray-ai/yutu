// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/playlist"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "playlist-list"
	listShort    = "List playlist's info"
	listLong     = "List playlist's info, such as title, description, etc"
	listIdsUsage = "Return the playlists with the given IDs for Stubby or Apiary"
	listCidUsage = "Return the playlists owned by the specified channel id"
)

type listIn struct {
	Ids                           []string `json:"ids"`
	ChannelId                     string   `json:"channelId"`
	Hl                            string   `json:"hl"`
	MaxResults                    int64    `json:"maxResults"`
	Mine                          *string  `json:"mine,omitempty"`
	OnBehalfOfContentOwner        string   `json:"onBehalfOfContentOwner"`
	OnBehalfOfContentOwnerChannel string   `json:"onBehalfOfContentOwnerChannel"`
	Parts                         []string `json:"parts"`
	Output                        string   `json:"output"`
	Jsonpath                      string   `json:"jsonpath"`
}

var listInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "channelId", "hl", "maxResults", "mine",
		"onBehalfOfContentOwner", "onBehalfOfContentOwnerChannel",
		"parts", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: listIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"channelId": {
			Type: "string", Description: listCidUsage,
			Default: json.RawMessage(`""`),
		},
		"hl": {
			Type: "string", Description: hlUsage,
			Default: json.RawMessage(`""`),
		},
		"maxResults": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: jsonschema.Ptr(float64(0)),
		},
		"mine": {
			Type: "string", Enum: []any{"true", "false", ""},
			Description: mineUsage, Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
		},
		"onBehalfOfContentOwnerChannel": {
			Type: "string", Description: "",
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
	playlistCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", listCidUsage)
	listCmd.Flags().StringVarP(&hl, "hl", "l", "", hlUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().BoolVarP(mine, "mine", "M", true, mineUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "", "",
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet", "status"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jpath, "jsonPath", "j", "", pkg.JPUsage)
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
	channelId = input.ChannelId
	hl = input.Hl
	maxResults = input.MaxResults
	mine = utils.BoolPtr(*input.Mine)
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	onBehalfOfContentOwnerChannel = input.OnBehalfOfContentOwnerChannel
	parts = input.Parts
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := list(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func list(writer io.Writer) error {
	p := playlist.NewPlaylist(
		playlist.WithIDs(ids),
		playlist.WithChannelId(channelId),
		playlist.WithHl(hl),
		playlist.WithMaxResults(maxResults),
		playlist.WithMine(mine),
		playlist.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		playlist.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
		playlist.WithService(nil),
	)

	return p.List(parts, output, jpath, writer)
}

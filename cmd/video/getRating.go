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
	getRatingTool  = "video-getRating"
	getRatingShort = "Get the rating of videos"
	getRatingLong  = "Get the rating of videos by ids"
	grIdsUsage     = "IDs of the videos to get the rating for"
)

type getRatingIn struct {
	Ids                    []string `json:"ids"`
	OnBehalfOfContentOwner string   `json:"onBehalfOfContentOwner"`
	Output                 string   `json:"output"`
	Jsonpath               string   `json:"jsonpath"`
}

var getRatingInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"ids", "onBehalfOfContentOwner", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Items: &jsonschema.Schema{
				Type: "string",
			},
			Description: grIdsUsage,
			Default:     json.RawMessage(`[]`),
		},
		"onBehalfOfContentOwner": {
			Type: "string", Description: "",
			Default: json.RawMessage(`""`),
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
			Name: getRatingTool, Title: getRatingShort, Description: getRatingLong,
			InputSchema: getRatingInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    true,
			},
		}, getRatingHandler,
	)
	videoCmd.AddCommand(getRatingCmd)

	getRatingCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, grIdsUsage)
	getRatingCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", "",
	)
	getRatingCmd.Flags().StringVarP(&output, "output", "o", "", pkg.TableUsage)
	getRatingCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = getRatingCmd.MarkFlagRequired("ids")
}

var getRatingCmd = &cobra.Command{
	Use:   "getRating",
	Short: getRatingShort,
	Long:  getRatingLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := getRating(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func getRatingHandler(
	ctx context.Context, req *mcp.CallToolRequest, input getRatingIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: getRatingTool, MinInterval: time.Second,
			},
		),
	)

	ids = input.Ids
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	output = input.Output
	jpath = input.Jsonpath

	var writer bytes.Buffer
	err := getRating(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func getRating(writer io.Writer) error {
	v := video.NewVideo(
		video.WithIDs(ids),
		video.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		video.WithService(nil),
	)

	return v.GetRating(output, jpath, writer)
}

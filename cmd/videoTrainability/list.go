// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoTrainability

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/videoTrainability"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool    = "videoTrainability-list"
	listShort   = "List video trainability"
	listLong    = "List video trainability. Use this tool to list video trainability by its ID."
	listExample = `  # List video trainability
  yutu videoTrainability list --videoId dQw4w9WgXcQ
  # List video trainability in JSON format
  yutu videoTrainability list --videoId dQw4w9WgXcQ --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"video_id"},
	Properties: map[string]*jsonschema.Schema{
		"video_id": {
			Type: "string", Description: "ID of the video",
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
			listTool, func(input videoTrainability.VideoTrainability, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	videoTrainabilityCmd.AddCommand(newListCmd())
}

type listOpts struct {
	videoId string
	output  string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.videoId, "videoId", "i", "", "ID of the video")
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
	_ = c.MarkFlagRequired("videoId")
}

func (o *listOpts) build() *videoTrainability.VideoTrainability {
	return videoTrainability.NewVideoTrainability(
		videoTrainability.WithVideoId(o.videoId),
		videoTrainability.WithOutput(o.output),
	).(*videoTrainability.VideoTrainability)
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

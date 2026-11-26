// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"encoding/json/jsontext"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "commentThread-list"
	listVidUsage = "Returns the comment threads of the specified video"
	listShort    = "List comment threads"
	listLong     = "List comment threads. Use this tool to list comment threads."
	listExample  = `  # List comment threads on a video
  yutu commentThread list --videoId dQw4w9WgXcQ --maxResults 10
  # Search comment threads by terms
  yutu commentThread list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --searchTerms 'hello'
  # List specific comment threads in JSON format
  yutu commentThread list --ids abc123,def456 --output json`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"all_threads_related_to_channel_id": {
			Type:        "string",
			Description: atrtcidUsage,
		},
		"channel_id": {Type: "string", Description: cidUsage},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: jsontext.Value("5"),
			Minimum: new(float64(0)),
		},
		"moderation_status": {
			Type:        "string",
			Enum:        []any{"published", "heldForReview", "likelySpam", "rejected"},
			Description: msUsage, Default: jsontext.Value(`"published"`),
		},
		"order": {
			Type: "string", Enum: []any{"orderUnspecified", "time", "relevance"},
			Description: orderUsage, Default: jsontext.Value(`"time"`),
		},
		"search_terms": {Type: "string", Description: stUsage},
		"text_format": {
			Type: "string", Enum: []any{"textFormatUnspecified", "html"},
			Description: tfUsage, Default: jsontext.Value(`"html"`),
		},
		"video_id": {Type: "string", Description: listVidUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: jsontext.Value(`["id","snippet"]`),
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
			listTool, func(input commentThread.CommentThread, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	commentThreadCmd.AddCommand(newListCmd())
}

type listOpts struct {
	ids                          []string
	allThreadsRelatedToChannelId string
	channelId                    string
	maxResults                   int64
	moderationStatus             string
	order                        string
	searchTerms                  string
	textFormat                   string
	videoId                      string
	parts                        []string
	output                       string
}

func (o *listOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, idsUsage)
	c.Flags().StringVarP(
		&o.allThreadsRelatedToChannelId, "allThreadsRelatedToChannelId", "a", "",
		atrtcidUsage,
	)
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().Int64VarP(&o.maxResults, "maxResults", "n", 5, pkg.MRUsage)
	c.Flags().StringVarP(
		&o.moderationStatus, "moderationStatus", "m", "published", msUsage,
	)
	c.Flags().StringVarP(&o.order, "order", "O", "time", orderUsage)
	c.Flags().StringVarP(&o.searchTerms, "searchTerms", "s", "", stUsage)
	c.Flags().StringVarP(&o.textFormat, "textFormat", "t", "html", tfUsage)
	c.Flags().StringVarP(&o.videoId, "videoId", "v", "", listVidUsage)
	c.Flags().StringSliceVarP(
		&o.parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "table", pkg.TableUsage)
}

func (o *listOpts) build() *commentThread.CommentThread {
	return commentThread.NewCommentThread(
		commentThread.WithIds(o.ids),
		commentThread.WithAllThreadsRelatedToChannelId(o.allThreadsRelatedToChannelId),
		commentThread.WithChannelId(o.channelId),
		commentThread.WithMaxResults(o.maxResults),
		commentThread.WithModerationStatus(o.moderationStatus),
		commentThread.WithOrder(o.order),
		commentThread.WithSearchTerms(o.searchTerms),
		commentThread.WithTextFormat(o.textFormat),
		commentThread.WithVideoId(o.videoId),
		commentThread.WithParts(o.parts),
		commentThread.WithOutput(o.output),
	).(*commentThread.CommentThread)
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

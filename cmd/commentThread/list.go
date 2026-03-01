// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"encoding/json"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/commentThread"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "commentThread-list"
	listShort    = "List YouTube comment threads"
	listLong     = "List YouTube comment threads\n\nExamples:\n  yutu commentThread list --videoId dQw4w9WgXcQ --maxResults 10\n  yutu commentThread list --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --searchTerms 'hello'\n  yutu commentThread list --ids abc123,def456 --output json"
	listVidUsage = "Returns the comment threads of the specified video"
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
			Default: json.RawMessage("5"),
			Minimum: new(float64(0)),
		},
		"moderation_status": {
			Type:        "string",
			Enum:        []any{"published", "heldForReview", "likelySpam", "rejected"},
			Description: msUsage, Default: json.RawMessage(`"published"`),
		},
		"order": {
			Type: "string", Enum: []any{"orderUnspecified", "time", "relevance"},
			Description: orderUsage, Default: json.RawMessage(`"time"`),
		},
		"search_terms": {Type: "string", Description: stUsage},
		"text_format": {
			Type: "string", Enum: []any{"textFormatUnspecified", "html"},
			Description: tfUsage, Default: json.RawMessage(`"html"`),
		},
		"video_id": {Type: "string", Description: listVidUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet"]`),
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
	commentThreadCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().StringVarP(
		&allThreadsRelatedToChannelId, "allThreadsRelatedToChannelId", "a", "",
		atrtcidUsage,
	)
	listCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	listCmd.Flags().Int64VarP(&maxResults, "maxResults", "n", 5, pkg.MRUsage)
	listCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "m", "published", msUsage,
	)
	listCmd.Flags().StringVarP(&order, "order", "O", "time", orderUsage)
	listCmd.Flags().StringVarP(&searchTerms, "searchTerms", "s", "", stUsage)
	listCmd.Flags().StringVarP(&textFormat, "textFormat", "t", "html", tfUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", listVidUsage)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: listShort,
	Long:  listLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := commentThread.NewCommentThread(
			commentThread.WithIds(ids),
			commentThread.WithAllThreadsRelatedToChannelId(allThreadsRelatedToChannelId),
			commentThread.WithChannelId(channelId),
			commentThread.WithMaxResults(maxResults),
			commentThread.WithModerationStatus(moderationStatus),
			commentThread.WithOrder(order),
			commentThread.WithSearchTerms(searchTerms),
			commentThread.WithTextFormat(textFormat),
			commentThread.WithVideoId(videoId),
			commentThread.WithParts(parts),
			commentThread.WithOutput(output),
			commentThread.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

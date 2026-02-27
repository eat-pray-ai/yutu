// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"encoding/json"
	"io"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	listTool     = "comment-list"
	listShort    = "List YouTube comments"
	listLong     = "List YouTube comments by ids\n\nExamples:\n  yutu comment list --parentId UgyXXXXXXXX --maxResults 10\n  yutu comment list --ids abc123,def456 --output json\n  yutu comment list --ids abc123 --textFormat plainText"
	listPidUsage = "Returns replies to the specified comment"
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"max_results": {
			Type: "number", Description: pkg.MRUsage,
			Default: json.RawMessage("5"),
			Minimum: new(float64(0)),
		},
		"parent_id": {Type: "string", Description: listPidUsage},
		"text_format": {
			Type: "string", Enum: []any{"textFormatUnspecified", "html", "plainText"},
			Description: tfUsage, Default: json.RawMessage(`"html"`),
		},
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
		}, cmd.GenToolHandler(
			listTool, func(input comment.Comment, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	commentCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	listCmd.Flags().Int64VarP(
		&maxResults, "maxResults", "n", 5, pkg.MRUsage,
	)
	listCmd.Flags().StringVarP(&parentId, "parentId", "P", "", listPidUsage)
	listCmd.Flags().StringVarP(
		&textFormat, "textFormat", "t", "html", tfUsage,
	)
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
		input := comment.NewComment(
			comment.WithIds(ids),
			comment.WithMaxResults(maxResults),
			comment.WithParentId(parentId),
			comment.WithTextFormat(textFormat),
			comment.WithParts(parts),
			comment.WithOutput(output),
			comment.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

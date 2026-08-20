// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/comment"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool    = "comment-update"
	updateIdUsage = "ID of the comment"
	updateShort   = "Update a comment on a video"
	updateLong    = "Update a comment on a video. Use this tool to update a comment on a video."
	updateExample = `# Update comment text
yutu comment update --id abc123 --textOriginal 'Updated comment'
# Like a comment
yutu comment update --id abc123 --viewerRating like
# Update comment text with rating enabled
yutu comment update --id abc123 --textOriginal 'New text' --canRate`
)

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: updateIdUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"can_rate":      {Type: "boolean", Description: crUsage},
		"text_original": {Type: "string", Description: toUsage},
		"viewer_rating": {
			Type: "string", Description: vrUsage,
			Enum: []any{"none", "like", "dislike"},
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent"},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: updateTool, Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			updateTool, func(input comment.Comment, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Update(writer)
			},
		),
	)
	commentCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringSliceVarP(&ids, "id", "i", []string{}, updateIdUsage)
	updateCmd.Flags().BoolVarP(canRate, "canRate", "R", false, crUsage)
	updateCmd.Flags().StringVarP(
		&textOriginal, "textOriginal", "t", "", toUsage,
	)
	updateCmd.Flags().StringVarP(
		&viewerRating, "viewerRating", "r", "", vrUsage,
	)
	updateCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = updateCmd.MarkFlagRequired("id")
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   updateShort,
	Long:    updateLong,
	Example: updateExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf("Would update comment: %s", strings.Join(ids, ", "))
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := comment.NewComment(
			comment.WithIds(ids),
			comment.WithCanRate(canRate),
			comment.WithTextOriginal(textOriginal),
			comment.WithViewerRating(viewerRating),
			comment.WithOutput(output),
		)
		utils.HandleCmdError(input.Update(c.OutOrStdout()), c)
	},
}

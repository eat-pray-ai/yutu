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
	smsTool    = "comment-setModerationStatus"
	smsShort   = "Set comment moderation status"
	smsLong    = "Set comment moderation status. Use this tool to set comment moderation status."
	smsExample = `# Publish a held comment
yutu comment setModerationStatus --ids abc123 --moderationStatus published
# Hold multiple comments for review
yutu comment setModerationStatus --ids abc123,def456 --moderationStatus heldForReview
# Reject a comment and ban author
yutu comment setModerationStatus --ids abc123 --moderationStatus rejected --banAuthor`
)

var setModerationStatusInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids", "moderation_status"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"moderation_status": {
			Type: "string", Description: msUsage,
			Enum: []any{"heldForReview", "published", "rejected"},
		},
		"ban_author": {Type: "boolean", Description: baUsage},
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
			Name: smsTool, Title: smsShort, Description: smsLong,
			InputSchema: setModerationStatusInSchema,
			Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			smsTool, func(input comment.Comment, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.SetModerationStatus(writer)
			},
		),
	)
	commentCmd.AddCommand(setModerationStatusCmd)

	setModerationStatusCmd.Flags().StringSliceVarP(
		&ids, "ids", "i", []string{}, idsUsage,
	)
	setModerationStatusCmd.Flags().StringVarP(
		&moderationStatus, "moderationStatus", "s", "", msUsage,
	)
	setModerationStatusCmd.Flags().BoolVarP(
		banAuthor, "banAuthor", "A", false, baUsage,
	)
	setModerationStatusCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	setModerationStatusCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = setModerationStatusCmd.MarkFlagRequired("ids")
	_ = setModerationStatusCmd.MarkFlagRequired("moderationStatus")
}

var setModerationStatusCmd = &cobra.Command{
	Use:     "setModerationStatus",
	Short:   smsShort,
	Long:    smsLong,
	Example: smsExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf(
			"Would set moderation status of comment(s): %s to %s",
			strings.Join(ids, ", "), moderationStatus,
		)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := comment.NewComment(
			comment.WithIds(ids),
			comment.WithModerationStatus(moderationStatus),
			comment.WithBanAuthor(banAuthor),
			comment.WithOutput(output),
		)
		utils.HandleCmdError(input.SetModerationStatus(c.OutOrStdout()), c)
	},
}

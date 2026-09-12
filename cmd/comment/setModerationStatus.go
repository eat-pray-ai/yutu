// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"encoding/json/jsontext"
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
	smsConfirm = "Set moderation status of comment(s): %s to %s"
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
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
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
		}, cobramcp.GenToolHandlerWithMRTR(
			smsTool, cobramcp.ConfirmThen(
				func(input comment.Comment) string {
					return fmt.Sprintf(
						smsConfirm, strings.Join(input.Ids, ", "), input.ModerationStatus,
					)
				},
				func(input comment.Comment, w io.Writer) error {
					return input.SetModerationStatus(w)
				},
			),
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
		return utils.ConfirmPreRun(
			c, fmt.Sprintf(smsConfirm, strings.Join(ids, ", "), moderationStatus),
		)
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

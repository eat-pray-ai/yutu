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
	"github.com/spf13/pflag"
)

const (
	smsTool    = "comment-setModerationStatus"
	smsConfirm = "Set moderation status of comment(s) %s to %s"
	smsShort   = "Set comment moderation status"
	smsLong    = "Set comment moderation status. Use this tool to set comment moderation status."
	smsExample = `  # Publish a held comment
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
	commentCmd.AddCommand(newSetModerationStatusCmd())
}

type setModerationStatusOpts struct {
	ids              []string
	moderationStatus string
	banAuthor        bool
	output           string
	yes              bool
}

func (o *setModerationStatusOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, idsUsage)
	c.Flags().StringVarP(
		&o.moderationStatus, "moderationStatus", "s", "", msUsage,
	)
	c.Flags().BoolVarP(&o.banAuthor, "banAuthor", "A", false, baUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
	_ = c.MarkFlagRequired("moderationStatus")
}

func (o *setModerationStatusOpts) build(f *pflag.FlagSet) *comment.Comment {
	return comment.NewComment(
		comment.WithIds(o.ids),
		comment.WithModerationStatus(o.moderationStatus),
		utils.OptIfChanged(f, "banAuthor", &o.banAuthor, comment.WithBanAuthor),
		comment.WithOutput(o.output),
	).(*comment.Comment)
}

func newSetModerationStatusCmd() *cobra.Command {
	var opts setModerationStatusOpts

	c := &cobra.Command{
		Use:     "setModerationStatus",
		Short:   smsShort,
		Long:    smsLong,
		Example: smsExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(
					smsConfirm, strings.Join(opts.ids, ", "), opts.moderationStatus,
				),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).SetModerationStatus(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

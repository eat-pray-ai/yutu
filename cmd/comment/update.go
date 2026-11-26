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
	updateTool    = "comment-update"
	updateIdUsage = "ID of the comment"
	updateConfirm = "Update comment: %s"
	updateShort   = "Update a comment"
	updateLong    = "Update a comment. Use this tool to update a comment."
	updateExample = `  # Update comment text
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
			Description: pkg.SilentUsage, Default: jsontext.Value(`"json"`),
		},
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
		}, cobramcp.GenToolHandlerWithMRTR(
			updateTool, cobramcp.ConfirmThen(
				func(input comment.Comment) string {
					return fmt.Sprintf(updateConfirm, strings.Join(input.Ids, ", "))
				},
				func(input comment.Comment, w io.Writer) error {
					comment.WithMaxResults(1)(&input)
					return input.Update(w)
				},
			),
		),
	)
	commentCmd.AddCommand(newUpdateCmd())
}

type updateOpts struct {
	ids          []string
	canRate      bool
	textOriginal string
	viewerRating string
	output       string
	yes          bool
}

func (o *updateOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "id", "i", nil, updateIdUsage)
	c.Flags().BoolVarP(&o.canRate, "canRate", "R", false, crUsage)
	c.Flags().StringVarP(&o.textOriginal, "textOriginal", "t", "", toUsage)
	c.Flags().StringVarP(&o.viewerRating, "viewerRating", "r", "", vrUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("id")
}

func (o *updateOpts) build(f *pflag.FlagSet) *comment.Comment {
	return comment.NewComment(
		comment.WithIds(o.ids),
		utils.OptIfChanged(f, "canRate", &o.canRate, comment.WithCanRate),
		comment.WithTextOriginal(o.textOriginal),
		comment.WithViewerRating(o.viewerRating),
		comment.WithMaxResults(1),
		comment.WithOutput(o.output),
	).(*comment.Comment)
}

func newUpdateCmd() *cobra.Command {
	var opts updateOpts

	c := &cobra.Command{
		Use:     "update",
		Short:   updateShort,
		Long:    updateLong,
		Example: updateExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(updateConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build(c.Flags()).Update(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

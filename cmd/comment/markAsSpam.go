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
	masTool    = "comment-markAsSpam"
	masConfirm = "Mark comment(s) %s as spam"
	masShort   = "Mark comments as spam"
	masLong    = "Mark comments as spam. Use this tool to mark comments as spam."
	masExample = `  # Mark a comment as spam
  yutu comment markAsSpam --ids abc123
  # Mark multiple comments as spam
  yutu comment markAsSpam --ids abc123,def456`
)

var markAsSpamInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: idsUsage,
			Items: &jsonschema.Schema{Type: "string"},
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
			Name: masTool, Title: masShort, Description: masLong,
			InputSchema: markAsSpamInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  false,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			masTool, cobramcp.ConfirmThen(
				func(input comment.Comment) string {
					return fmt.Sprintf(masConfirm, strings.Join(input.Ids, ", "))
				},
				func(input comment.Comment, w io.Writer) error {
					return input.MarkAsSpam(w)
				},
			),
		),
	)
	commentCmd.AddCommand(newMarkAsSpamCmd())
}

type markAsSpamOpts struct {
	ids    []string
	output string
	yes    bool
}

func (o *markAsSpamOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, idsUsage)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *markAsSpamOpts) build() *comment.Comment {
	return comment.NewComment(
		comment.WithIds(o.ids),
		comment.WithOutput(o.output),
	).(*comment.Comment)
}

func newMarkAsSpamCmd() *cobra.Command {
	var opts markAsSpamOpts

	c := &cobra.Command{
		Use:     "markAsSpam",
		Short:   masShort,
		Long:    masLong,
		Example: masExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(masConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().MarkAsSpam(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

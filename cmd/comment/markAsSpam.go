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
	masTool    = "comment-markAsSpam"
	masShort   = "Mark comments as spam"
	masLong    = "Mark comments as spam. Use this tool to mark comments as spam."
	masExample = `# Mark a comment as spam
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
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"confirmed": {Type: "boolean", Description: pkg.ConfirmedUsage},
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
		}, cobramcp.GenToolHandler(
			masTool, func(input comment.Comment, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.MarkAsSpam(writer)
			},
		),
	)
	commentCmd.AddCommand(markAsSpamCmd)

	markAsSpamCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, idsUsage)
	markAsSpamCmd.Flags().StringP("output", "o", "", pkg.SilentUsage)
	markAsSpamCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = markAsSpamCmd.MarkFlagRequired("ids")
}

var markAsSpamCmd = &cobra.Command{
	Use:     "markAsSpam",
	Short:   masShort,
	Long:    masLong,
	Example: masExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf(
			"Would mark comment(s) as spam: %s", strings.Join(ids, ", "),
		)
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		output, _ := c.Flags().GetString("output")
		input := comment.NewComment(
			comment.WithIds(ids),
			comment.WithOutput(output),
		)
		utils.HandleCmdError(input.MarkAsSpam(c.OutOrStdout()), c)
	},
}

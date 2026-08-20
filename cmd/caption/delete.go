// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"fmt"
	"io"
	"strings"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	deleteTool     = "caption-delete"
	deleteIdsUsage = "IDs of the captions to delete"
	deleteShort    = "Delete captions"
	deleteLong     = "Delete captions. Use this tool to delete captions of a video by IDs."
	deleteExample  = `# Delete a caption by ID
yutu caption delete --ids abc123
# Delete multiple captions
yutu caption delete --ids abc123,def456`
)

var deleteInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"ids"},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: deleteIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"confirmed":                  {Type: "boolean", Description: pkg.ConfirmedUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: deleteTool, Title: deleteShort, Description: deleteLong,
			InputSchema: deleteInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(true),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandler(
			deleteTool, func(input caption.Caption, writer io.Writer) error {
				if !input.Confirmed {
					return utils.ErrNotConfirmed
				}
				return input.Delete(writer)
			},
		),
	)
	captionCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, deleteIdsUsage)
	deleteCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	deleteCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	deleteCmd.Flags().Bool("yes", false, pkg.ConfirmedUsage)
	_ = deleteCmd.MarkFlagRequired("ids")
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   deleteShort,
	Long:    deleteLong,
	Example: deleteExample,
	PreRunE: func(c *cobra.Command, _ []string) error {
		msg := fmt.Sprintf("Would delete caption(s): %s", strings.Join(ids, ", "))
		return utils.ConfirmPreRun(c, msg)
	},
	Run: func(c *cobra.Command, _ []string) {
		input := caption.NewCaption(
			caption.WithIds(ids),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		)
		utils.HandleCmdError(input.Delete(c.OutOrStdout()), c)
	},
}

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
	deleteConfirm  = "Delete caption(s): %s"
	deleteShort    = "Delete captions"
	deleteLong     = "Delete captions. Use this tool to delete captions of a video by IDs."
	deleteExample  = `  # Delete a caption by ID
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
		}, cobramcp.GenToolHandlerWithMRTR(
			deleteTool, cobramcp.ConfirmThen(
				func(input caption.Caption) string {
					return fmt.Sprintf(deleteConfirm, strings.Join(input.Ids, ", "))
				},
				func(input caption.Caption, w io.Writer) error {
					return input.Delete(w)
				},
			),
		),
	)
	captionCmd.AddCommand(newDeleteCmd())
}

type deleteOpts struct {
	ids                    []string
	onBehalfOf             string
	onBehalfOfContentOwner string
	yes                    bool
}

func (o *deleteOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringSliceVarP(&o.ids, "ids", "i", nil, deleteIdsUsage)
	c.Flags().StringVarP(&o.onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	c.Flags().StringVarP(
		&o.onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("ids")
}

func (o *deleteOpts) build() *caption.Caption {
	return caption.NewCaption(
		caption.WithIds(o.ids),
		caption.WithOnBehalfOf(o.onBehalfOf),
		caption.WithOnBehalfOfContentOwner(o.onBehalfOfContentOwner),
	).(*caption.Caption)
}

func newDeleteCmd() *cobra.Command {
	var opts deleteOpts

	c := &cobra.Command{
		Use:     "delete",
		Short:   deleteShort,
		Long:    deleteLong,
		Example: deleteExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(deleteConfirm, strings.Join(opts.ids, ", ")),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Delete(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

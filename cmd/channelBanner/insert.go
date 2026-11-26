// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"encoding/json/jsontext"
	"fmt"
	"io"

	cobramcp "github.com/eat-pray-ai/cobra-mcp"
	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/channelBanner"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool    = "channelBanner-insert"
	insertConfirm = "Upload channel banner %s for channel %s"
	insertShort   = "Upload a channel banner"
	insertLong    = "Upload a channel banner. Use this tool to upload a channel banner."
	insertExample = `  # Upload a channel banner
  yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.jpg
  # Upload a channel banner with JSON output
  yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.png --output json`
)

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"channel_id", "file"},
	Properties: map[string]*jsonschema.Schema{
		"channel_id": {Type: "string", Description: cidUsage},
		"file":       {Type: "string", Description: fileUsage},

		"on_behalf_of_content_owner": {
			Type:        "string",
			Description: pkg.OBOCOUsage,
		},
		"on_behalf_of_content_owner_channel": {
			Type:        "string",
			Description: pkg.OBOCOCUsage,
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
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    false,
			},
		}, cobramcp.GenToolHandlerWithMRTR(
			insertTool, cobramcp.ConfirmThen(
				func(input channelBanner.ChannelBanner) string {
					return fmt.Sprintf(insertConfirm, input.File, input.ChannelId)
				},
				func(input channelBanner.ChannelBanner, w io.Writer) error {
					return input.Insert(w)
				},
			),
		),
	)
	channelBannerCmd.AddCommand(newInsertCmd())
}

type insertOpts struct {
	channelId string
	file      string
	oboco     string
	obococ    string
	output    string
	yes       bool
}

func (o *insertOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(&o.channelId, "channelId", "c", "", cidUsage)
	c.Flags().StringVarP(&o.file, "file", "f", "", fileUsage)
	c.Flags().StringVarP(
		&o.oboco, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	c.Flags().StringVarP(
		&o.obococ, "onBehalfOfContentOwnerChannel", "B", "",
		pkg.OBOCOCUsage,
	)
	c.Flags().StringVarP(&o.output, "output", "o", "", pkg.SilentUsage)
	c.Flags().BoolVar(&o.yes, "yes", false, pkg.ConfirmedUsage)
	_ = c.MarkFlagRequired("channelId")
	_ = c.MarkFlagRequired("file")
}

func (o *insertOpts) build() *channelBanner.ChannelBanner {
	return channelBanner.NewChannelBanner(
		channelBanner.WithChannelId(o.channelId),
		channelBanner.WithFile(o.file),
		channelBanner.WithOBOCO(o.oboco),
		channelBanner.WithOBOCOC(o.obococ),
		channelBanner.WithOutput(o.output),
	).(*channelBanner.ChannelBanner)
}

func newInsertCmd() *cobra.Command {
	var opts insertOpts

	c := &cobra.Command{
		Use:     "insert",
		Short:   insertShort,
		Long:    insertLong,
		Example: insertExample,
		PreRunE: func(c *cobra.Command, _ []string) error {
			return utils.ConfirmPreRun(
				c, opts.yes, fmt.Sprintf(insertConfirm, opts.file, opts.channelId),
			)
		},
		RunE: func(c *cobra.Command, _ []string) error {
			return opts.build().Insert(c.OutOrStdout())
		},
	}

	opts.bindFlags(c)
	return c
}

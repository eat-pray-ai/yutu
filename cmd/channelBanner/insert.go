// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"encoding/json"
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
	insertShort   = "Insert a YouTube channel banner"
	insertLong    = "Insert a YouTube channel banner. Use this tool when you need to insert or upload a channel banner."
	insertExample = `yutu channelBanner insert --channelId UC_x5XG1OV2P6uZZ5FSM9Ttw --file banner.jpg
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
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
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
		}, cobramcp.GenToolHandler(
			insertTool,
			func(input channelBanner.ChannelBanner, writer io.Writer) error {
				return input.Insert(writer)
			},
		),
	)
	channelBannerCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&channelId, "channelId", "c", "", cidUsage)
	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "b", "", pkg.OBOCOUsage,
	)
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwnerChannel, "onBehalfOfContentOwnerChannel", "B", "",
		pkg.OBOCOCUsage,
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("channelId")
	_ = insertCmd.MarkFlagRequired("file")
}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Short:   insertShort,
	Long:    insertLong,
	Example: insertExample,
	Run: func(cmd *cobra.Command, args []string) {
		input := channelBanner.NewChannelBanner(
			channelBanner.WithChannelId(channelId),
			channelBanner.WithFile(file),
			channelBanner.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			channelBanner.WithOnBehalfOfContentOwnerChannel(onBehalfOfContentOwnerChannel),
			channelBanner.WithOutput(output),
			channelBanner.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.Insert(cmd.OutOrStdout()), cmd)
	},
}

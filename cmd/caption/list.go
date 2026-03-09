// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"encoding/json"
	"io"

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
	listTool     = "caption-list"
	listIdsUsage = "IDs of the captions to list"
	listShort    = "List captions"
	listLong     = "List captions. Use this tool to list captions of a video."
	listExample  = `# List captions of a video
yutu caption list --videoId dQw4w9WgXcQ
# List captions in JSON format
yutu caption list --videoId dQw4w9WgXcQ --output json
# List specific captions by IDs
yutu caption list --ids abc123,def456 --videoId dQw4w9WgXcQ`
)

var listInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{},
	Properties: map[string]*jsonschema.Schema{
		"ids": {
			Type: "array", Description: listIdsUsage,
			Items: &jsonschema.Schema{Type: "string"},
		},
		"video_id":                   {Type: "string", Description: vidUsage},
		"on_behalf_of":               {Type: "string", Description: pkg.OBOUsage},
		"on_behalf_of_content_owner": {Type: "string", Description: pkg.OBOCOUsage},
		"parts": {
			Type: "array", Description: pkg.PartsUsage,
			Items:   &jsonschema.Schema{Type: "string"},
			Default: json.RawMessage(`["id","snippet"]`),
		},
		"output": {
			Type: "string", Description: pkg.TableUsage,
			Enum:    []any{"json", "yaml", "table"},
			Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: listTool, Title: listShort, Description: listLong,
			InputSchema: listInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: new(false),
				IdempotentHint:  true,
				OpenWorldHint:   new(true),
				ReadOnlyHint:    true,
			},
		}, cobramcp.GenToolHandler(
			listTool, func(input caption.Caption, writer io.Writer) error {
				return input.List(writer)
			},
		),
	)
	captionCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceVarP(&ids, "ids", "i", []string{}, listIdsUsage)
	listCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	listCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", pkg.OBOUsage)
	listCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", pkg.OBOCOUsage,
	)
	listCmd.Flags().StringSliceVarP(
		&parts, "parts", "p", []string{"id", "snippet"}, pkg.PartsUsage,
	)
	listCmd.Flags().StringVarP(&output, "output", "o", "table", pkg.TableUsage)
	listCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   listShort,
	Long:    listLong,
	Example: listExample,
	Run: func(cmd *cobra.Command, args []string) {
		input := caption.NewCaption(
			caption.WithIds(ids),
			caption.WithVideoId(videoId),
			caption.WithOnBehalfOf(onBehalfOf),
			caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
			caption.WithParts(parts),
			caption.WithOutput(output),
			caption.WithJsonpath(jsonpath),
		)
		utils.HandleCmdError(input.List(cmd.OutOrStdout()), cmd)
	},
}

package thumbnail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/thumbnail"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type setIn struct {
	File     string `json:"file"`
	VideoId  string `json:"videoId"`
	Output   string `json:"output"`
	Jsonpath string `json:"jsonpath"`
}

var setInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"file", "videoId", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"file": {
			Type: "string", Description: fileUsage,
			Default: json.RawMessage(`""`),
		},
		"videoId": {
			Type: "string", Description: vidUsage,
			Default: json.RawMessage(`""`),
		},
		"output": {
			Type: "string", Enum: []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage, Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type: "string", Description: pkg.JPUsage,
			Default: json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "thumbnail-set", Title: short, Description: long,
			InputSchema: setInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  true,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, setHandler,
	)
	thumbnailCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	setCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	setCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	setCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = setCmd.MarkFlagRequired("file")
	_ = setCmd.MarkFlagRequired("videoId")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: short,
	Long:  long,
	Run: func(cmd *cobra.Command, args []string) {
		err := set(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func setHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input setIn,
) (*mcp.CallToolResult, any, error) {
	file = input.File
	videoId = input.VideoId
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "thumbnail set started")

	var writer bytes.Buffer
	err := set(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "thumbnail set failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "thumbnail set completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func set(writer io.Writer) error {
	t := thumbnail.NewThumbnail(
		thumbnail.WithFile(file),
		thumbnail.WithVideoId(videoId),
		thumbnail.WithService(nil),
	)

	return t.Set(output, jpath, writer)
}

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateShort = "Update caption"
	updateLong  = "Update caption of a video"
)

type updateIn struct {
	File                   string  `json:"file"`
	AudioTrackType         string  `json:"audioTrackType"`
	IsAutoSynced           *string `json:"isAutoSynced,omitempty"`
	IsCC                   *string `json:"isCC,omitempty"`
	IsDraft                *string `json:"isDraft,omitempty"`
	IsEasyReader           *string `json:"isEasyReader,omitempty"`
	IsLarge                *string `json:"isLarge,omitempty"`
	Language               string  `json:"language"`
	Name                   string  `json:"name"`
	TrackKind              string  `json:"trackKind"`
	VideoId                string  `json:"videoId"`
	OnBehalfOf             string  `json:"onBehalfOf"`
	OnBehalfOfContentOwner string  `json:"onBehalfOfContentOwner"`
	Output                 string  `json:"output"`
	Jsonpath               string  `json:"jsonpath"`
}

var updateInSchema = &jsonschema.Schema{
	Type: "object",
	Required: []string{
		"file", "audioTrackType", "isAutoSynced", "isCC", "isDraft",
		"isEasyReader", "isLarge", "language", "name", "trackKind",
		"videoId", "onBehalfOf", "onBehalfOfContentOwner", "output", "jsonpath",
	},
	Properties: map[string]*jsonschema.Schema{
		"file": {
			Type:        "string",
			Description: fileUsage,
			Default:     json.RawMessage(`""`),
		},
		"audioTrackType": {
			Type:        "string",
			Description: attUsage,
			Default:     json.RawMessage(`"unknown"`),
		},
		"isAutoSynced": {
			Type:        "string",
			Enum:        []any{"true", "false", ""},
			Description: iasUsage,
			Default:     json.RawMessage(`""`),
		},
		"isCC": {
			Type:        "string",
			Enum:        []any{"true", "false", ""},
			Description: iscUsage,
			Default:     json.RawMessage(`""`),
		},
		"isDraft": {
			Type:        "string",
			Enum:        []any{"true", "false", ""},
			Description: isdUsage,
			Default:     json.RawMessage(`""`),
		},
		"isEasyReader": {
			Type:        "string",
			Enum:        []any{"true", "false", ""},
			Description: iserUsage,
			Default:     json.RawMessage(`""`),
		},
		"isLarge": {
			Type:        "string",
			Enum:        []any{"true", "false", ""},
			Description: islUsage,
			Default:     json.RawMessage(`""`),
		},
		"language": {
			Type:        "string",
			Description: languageUsage,
			Default:     json.RawMessage(`""`),
		},
		"name": {
			Type:        "string",
			Description: nameUsage,
			Default:     json.RawMessage(`""`),
		},
		"trackKind": {
			Type:        "string",
			Enum:        []any{"standard", "ASR", "forced"},
			Description: tkUsage,
			Default:     json.RawMessage(`"standard"`),
		},
		"videoId": {
			Type:        "string",
			Description: vidUsage,
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOf": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
		"onBehalfOfContentOwner": {
			Type:        "string",
			Description: "",
			Default:     json.RawMessage(`""`),
		},
		"output": {
			Type:        "string",
			Enum:        []any{"json", "yaml", "silent", ""},
			Description: pkg.SilentUsage,
			Default:     json.RawMessage(`"yaml"`),
		},
		"jsonpath": {
			Type:        "string",
			Description: pkg.JPUsage,
			Default:     json.RawMessage(`""`),
		},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: "caption-update", Title: updateShort, Description: updateLong,
			InputSchema: updateInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, updateHandler,
	)
	captionCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	updateCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	updateCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	updateCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	updateCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	updateCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	updateCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	updateCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	updateCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	updateCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	updateCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	updateCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	updateCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	updateCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	updateCmd.Flags().StringVarP(&jpath, "jsonpath", "j", "", pkg.JPUsage)

	_ = updateCmd.MarkFlagRequired("videoId")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: updateShort,
	Long:  updateLong,
	Run: func(cmd *cobra.Command, args []string) {
		err := update(cmd.OutOrStdout())
		if err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func updateHandler(
	ctx context.Context, _ *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	file = input.File
	audioTrackType = input.AudioTrackType
	isAutoSynced = utils.BoolPtr(*input.IsAutoSynced)
	isCC = utils.BoolPtr(*input.IsCC)
	isDraft = utils.BoolPtr(*input.IsDraft)
	isEasyReader = utils.BoolPtr(*input.IsEasyReader)
	isLarge = utils.BoolPtr(*input.IsLarge)
	language = input.Language
	name = input.Name
	trackKind = input.TrackKind
	videoId = input.VideoId
	onBehalfOf = input.OnBehalfOf
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	output = input.Output
	jpath = input.Jsonpath

	slog.InfoContext(ctx, "caption update started")

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		slog.ErrorContext(
			ctx, "caption update failed", "error", err, "input", input,
		)
		return nil, nil, err
	}
	slog.InfoContext(
		ctx, "caption update completed successfully",
		"resultSize", writer.Len(),
	)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func update(writer io.Writer) error {
	c := caption.NewCation(
		caption.WithFile(file),
		caption.WithAudioTrackType(audioTrackType),
		caption.WithIsAutoSynced(isAutoSynced),
		caption.WithIsCC(isCC),
		caption.WithIsDraft(isDraft),
		caption.WithIsEasyReader(isEasyReader),
		caption.WithIsLarge(isLarge),
		caption.WithLanguage(language),
		caption.WithName(name),
		caption.WithTrackKind(trackKind),
		caption.WithOnBehalfOf(onBehalfOf),
		caption.WithOnBehalfOfContentOwner(onBehalfOfContentOwner),
		caption.WithVideoId(videoId),
		caption.WithService(nil),
	)

	return c.Update(output, jpath, writer)
}

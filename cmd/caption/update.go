// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/eat-pray-ai/yutu/cmd"
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/caption"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	updateTool  = "caption-update"
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
	Type:     "object",
	Required: []string{"videoId"},
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
			Name: updateTool, Title: updateShort, Description: updateLong,
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
	updateCmd.Flags().StringVarP(&jsonpath, "jsonpath", "j", "", pkg.JPUsage)

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
	ctx context.Context, req *mcp.CallToolRequest, input updateIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: updateTool, MinInterval: time.Second,
			},
		),
	)

	file = input.File
	audioTrackType = input.AudioTrackType
	isAutoSynced = utils.StrToBoolPtr(input.IsAutoSynced)
	isCC = utils.StrToBoolPtr(input.IsCC)
	isDraft = utils.StrToBoolPtr(input.IsDraft)
	isEasyReader = utils.StrToBoolPtr(input.IsEasyReader)
	isLarge = utils.StrToBoolPtr(input.IsLarge)
	language = input.Language
	name = input.Name
	trackKind = input.TrackKind
	videoId = input.VideoId
	onBehalfOf = input.OnBehalfOf
	onBehalfOfContentOwner = input.OnBehalfOfContentOwner
	output = input.Output
	jsonpath = input.Jsonpath

	var writer bytes.Buffer
	err := update(&writer)
	if err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
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

	return c.Update(output, jsonpath, writer)
}

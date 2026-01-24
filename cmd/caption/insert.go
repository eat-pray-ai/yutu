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
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

const (
	insertTool  = "caption-insert"
	insertShort = "Insert caption"
	insertLong  = "Insert caption to a video"
)

type insertIn struct {
	File                   string `json:"file"`
	AudioTrackType         string `json:"audioTrackType"`
	IsAutoSynced           *bool  `json:"isAutoSynced,omitempty"`
	IsCC                   *bool  `json:"isCC,omitempty"`
	IsDraft                *bool  `json:"isDraft,omitempty"`
	IsEasyReader           *bool  `json:"isEasyReader,omitempty"`
	IsLarge                *bool  `json:"isLarge,omitempty"`
	Language               string `json:"language"`
	Name                   string `json:"name"`
	TrackKind              string `json:"trackKind"`
	VideoId                string `json:"videoId"`
	OnBehalfOf             string `json:"onBehalfOf"`
	OnBehalfOfContentOwner string `json:"onBehalfOfContentOwner"`
	Output                 string `json:"output"`
	Jsonpath               string `json:"jsonpath"`
}

var insertInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"file", "videoId"},
	Properties: map[string]*jsonschema.Schema{
		"file": {Type: "string", Description: fileUsage},
		"audioTrackType": {
			Type: "string", Description: attUsage,
			Enum:    []any{"unknown", "primary", "commentary", "descriptive"},
			Default: json.RawMessage(`"unknown"`),
		},
		"isAutoSynced": {
			Type: "boolean", Description: iasUsage,
			Enum: []any{true, false},
		},
		"isCC": {
			Type: "boolean", Description: iscUsage,
			Enum: []any{true, false},
		},
		"isDraft": {
			Type: "boolean", Description: isdUsage,
			Enum: []any{true, false},
		},
		"isEasyReader": {
			Type: "boolean", Description: iserUsage,
			Enum: []any{true, false},
		},
		"isLarge": {
			Type: "boolean", Description: islUsage,
			Enum: []any{true, false},
		},
		"language": {Type: "string", Description: languageUsage},
		"name":     {Type: "string", Description: nameUsage},
		"trackKind": {
			Type: "string", Description: tkUsage,
			Enum:    []any{"standard", "ASR", "forced"},
			Default: json.RawMessage(`"standard"`),
		},
		"videoId": {
			Type: "string", Description: vidUsage,
			Default: json.RawMessage(`""`),
		},
		"onBehalfOf":             {Type: "string"},
		"onBehalfOfContentOwner": {Type: "string"},
		"output": {
			Type: "string", Description: pkg.SilentUsage,
			Enum:    []any{"json", "yaml", "silent", ""},
			Default: json.RawMessage(`"yaml"`),
		},
		"jsonpath": {Type: "string", Description: pkg.JPUsage},
	},
}

func init() {
	mcp.AddTool(
		cmd.Server, &mcp.Tool{
			Name: insertTool, Title: insertShort, Description: insertLong,
			InputSchema: insertInSchema, Annotations: &mcp.ToolAnnotations{
				DestructiveHint: jsonschema.Ptr(false),
				IdempotentHint:  false,
				OpenWorldHint:   jsonschema.Ptr(true),
				ReadOnlyHint:    false,
			},
		}, insertHandler,
	)
	captionCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringVarP(&file, "file", "f", "", fileUsage)
	insertCmd.Flags().StringVarP(
		&audioTrackType, "audioTrackType", "a", "unknown", attUsage,
	)
	insertCmd.Flags().BoolVarP(
		isAutoSynced, "isAutoSynced", "A", true, iasUsage,
	)
	insertCmd.Flags().BoolVarP(isCC, "isCC", "C", false, iscUsage)
	insertCmd.Flags().BoolVarP(isDraft, "isDraft", "D", false, isdUsage)
	insertCmd.Flags().BoolVarP(
		isEasyReader, "isEasyReader", "E", false, iserUsage,
	)
	insertCmd.Flags().BoolVarP(isLarge, "isLarge", "L", false, islUsage)
	insertCmd.Flags().StringVarP(&language, "language", "l", "", languageUsage)
	insertCmd.Flags().StringVarP(&name, "name", "n", "", nameUsage)
	insertCmd.Flags().StringVarP(
		&trackKind, "trackKind", "t", "standard", tkUsage,
	)
	insertCmd.Flags().StringVarP(&videoId, "videoId", "v", "", vidUsage)
	insertCmd.Flags().StringVarP(&onBehalfOf, "onBehalfOf", "b", "", "")
	insertCmd.Flags().StringVarP(
		&onBehalfOfContentOwner, "onBehalfOfContentOwner", "B", "", "",
	)
	insertCmd.Flags().StringVarP(&output, "output", "o", "", pkg.SilentUsage)
	insertCmd.Flags().StringVarP(&jsonpath, "jsonPath", "j", "", pkg.JPUsage)

	_ = insertCmd.MarkFlagRequired("file")
	_ = insertCmd.MarkFlagRequired("videoId")
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: insertShort,
	Long:  insertLong,
	Run: func(cmd *cobra.Command, args []string) {
		input := &insertIn{
			File:                   file,
			AudioTrackType:         audioTrackType,
			IsAutoSynced:           isAutoSynced,
			IsCC:                   isCC,
			IsDraft:                isDraft,
			IsEasyReader:           isEasyReader,
			IsLarge:                isLarge,
			Language:               language,
			Name:                   name,
			TrackKind:              trackKind,
			VideoId:                videoId,
			OnBehalfOf:             onBehalfOf,
			OnBehalfOfContentOwner: onBehalfOfContentOwner,
			Output:                 output,
			Jsonpath:               jsonpath,
		}
		if err := input.call(cmd.OutOrStdout()); err != nil {
			_ = cmd.Help()
			cmd.PrintErrf("Error: %v\n", err)
		}
	},
}

func insertHandler(
	ctx context.Context, req *mcp.CallToolRequest, input insertIn,
) (*mcp.CallToolResult, any, error) {
	logger := slog.New(
		mcp.NewLoggingHandler(
			req.Session,
			&mcp.LoggingHandlerOptions{
				LoggerName: insertTool, MinInterval: time.Second,
			},
		),
	)

	var writer bytes.Buffer
	if err := input.call(&writer); err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func (i *insertIn) call(writer io.Writer, opts ...caption.Option) error {
	defaultOpts := []caption.Option{
		caption.WithFile(i.File),
		caption.WithAudioTrackType(i.AudioTrackType),
		caption.WithIsAutoSynced(i.IsAutoSynced),
		caption.WithIsCC(i.IsCC),
		caption.WithIsDraft(i.IsDraft),
		caption.WithIsEasyReader(i.IsEasyReader),
		caption.WithIsLarge(i.IsLarge),
		caption.WithLanguage(i.Language),
		caption.WithName(i.Name),
		caption.WithTrackKind(i.TrackKind),
		caption.WithOnBehalfOf(i.OnBehalfOf),
		caption.WithOnBehalfOfContentOwner(i.OnBehalfOfContentOwner),
		caption.WithVideoId(i.VideoId),
		caption.WithService(nil),
	}
	defaultOpts = append(defaultOpts, opts...)

	c := caption.NewCation(defaultOpts...)

	return c.Insert(i.Output, i.Jsonpath, writer)
}

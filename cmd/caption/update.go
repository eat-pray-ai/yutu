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
	updateTool  = "caption-update"
	updateShort = "Update caption"
	updateLong  = "Update caption of a video"
)

type updateIn struct {
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

var updateInSchema = &jsonschema.Schema{
	Type:     "object",
	Required: []string{"videoId"},
	Properties: map[string]*jsonschema.Schema{
		"file": {Type: "string", Description: fileUsage},
		"audioTrackType": {
			Type: "string", Description: attUsage,
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
		"videoId":                {Type: "string", Description: vidUsage},
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
		input := &updateIn{
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

	var writer bytes.Buffer
	if err := input.call(&writer); err != nil {
		logger.ErrorContext(ctx, err.Error(), "input", input)
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: writer.String()}}}, nil, nil
}

func (u *updateIn) call(writer io.Writer, opts ...caption.Option) error {
	defaultOpts := []caption.Option{
		caption.WithFile(u.File),
		caption.WithAudioTrackType(u.AudioTrackType),
		caption.WithIsAutoSynced(u.IsAutoSynced),
		caption.WithIsCC(u.IsCC),
		caption.WithIsDraft(u.IsDraft),
		caption.WithIsEasyReader(u.IsEasyReader),
		caption.WithIsLarge(u.IsLarge),
		caption.WithLanguage(u.Language),
		caption.WithName(u.Name),
		caption.WithTrackKind(u.TrackKind),
		caption.WithOnBehalfOf(u.OnBehalfOf),
		caption.WithOnBehalfOfContentOwner(u.OnBehalfOfContentOwner),
		caption.WithVideoId(u.VideoId),
		caption.WithService(nil),
	}
	defaultOpts = append(defaultOpts, opts...)

	c := caption.NewCation(defaultOpts...)

	return c.Update(u.Output, u.Jsonpath, writer)
}

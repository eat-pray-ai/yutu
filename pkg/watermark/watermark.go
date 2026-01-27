// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
)

var (
	errSetWatermark   = errors.New("failed to set watermark")
	errUnsetWatermark = errors.New("failed to unset watermark")
)

type Watermark struct {
	*pkg.DefaultFields
	ChannelId              string `yaml:"channel_id" json:"channel_id"`
	File                   string `yaml:"file" json:"file"`
	InVideoPosition        string `yaml:"in_video_position" json:"in_video_position"`
	DurationMs             uint64 `yaml:"duration_ms" json:"duration_ms"`
	OffsetMs               uint64 `yaml:"offset_ms" json:"offset_ms"`
	OffsetType             string `yaml:"offset_type" json:"offset_type"`
	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
}

type IWatermark interface {
	Set(io.Writer) error
	Unset(io.Writer) error
	GetDefaultFields() *pkg.DefaultFields
	preRun()
}

type Option func(*Watermark)

func NewWatermark(opts ...Option) IWatermark {
	w := &Watermark{DefaultFields: &pkg.DefaultFields{}}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Watermark) GetDefaultFields() *pkg.DefaultFields {
	return w.DefaultFields
}

func (w *Watermark) preRun() {
	if w.Service == nil {
		w.Service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (w *Watermark) Set(writer io.Writer) error {
	w.preRun()
	file, err := pkg.Root.Open(w.File)
	if err != nil {
		return errors.Join(errSetWatermark, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	inVideoBranding := &youtube.InvideoBranding{
		Position: &youtube.InvideoPosition{},
		Timing:   &youtube.InvideoTiming{},
	}
	if w.InVideoPosition != "" {
		inVideoBranding.Position.Type = "corner"
		inVideoBranding.Position.CornerPosition = w.InVideoPosition
	}
	if w.DurationMs != 0 {
		inVideoBranding.Timing.DurationMs = w.DurationMs
	}
	if w.OffsetMs != 0 {
		inVideoBranding.Timing.OffsetMs = w.OffsetMs
	}
	if w.OffsetType != "" {
		inVideoBranding.Timing.Type = w.OffsetType
	}

	call := w.Service.Watermarks.Set(w.ChannelId, inVideoBranding).Media(file)
	if w.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.OnBehalfOfContentOwner)
	}

	err = call.Do()
	if err != nil {
		return errors.Join(errSetWatermark, err)
	}

	_, _ = fmt.Fprintf(writer, "Watermark set for channel %s\n", w.ChannelId)
	return nil
}

func (w *Watermark) Unset(writer io.Writer) error {
	w.preRun()
	call := w.Service.Watermarks.Unset(w.ChannelId)
	if w.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.OnBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		return errors.Join(errUnsetWatermark, err)
	}

	_, _ = fmt.Fprintf(writer, "Watermark unset for channel %s\n", w.ChannelId)
	return nil
}

func WithChannelId(channelId string) Option {
	return func(w *Watermark) {
		w.ChannelId = channelId
	}
}

func WithFile(file string) Option {
	return func(w *Watermark) {
		w.File = file
	}
}

func WithInVideoPosition(inVideoPosition string) Option {
	return func(w *Watermark) {
		w.InVideoPosition = inVideoPosition
	}
}

func WithDurationMs(durationMs uint64) Option {
	return func(w *Watermark) {
		w.DurationMs = durationMs
	}
}

func WithOffsetMs(offsetMs uint64) Option {
	return func(w *Watermark) {
		w.OffsetMs = offsetMs
	}
}

func WithOffsetType(offsetType string) Option {
	return func(w *Watermark) {
		w.OffsetType = offsetType
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(w *Watermark) {
		w.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

var WithService = pkg.WithService[*Watermark]

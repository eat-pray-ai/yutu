package watermark

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service           *youtube.Service
	errSetWatermark   = errors.New("failed to set watermark")
	errUnsetWatermark = errors.New("failed to unset watermark")
)

type watermark struct {
	ChannelId              string `yaml:"channel_id" json:"channel_id"`
	File                   string `yaml:"file" json:"file"`
	InVideoPosition        string `yaml:"in_video_position" json:"in_video_position"`
	DurationMs             uint64 `yaml:"duration_ms" json:"duration_ms"`
	OffsetMs               uint64 `yaml:"offset_ms" json:"offset_ms"`
	OffsetType             string `yaml:"offset_type" json:"offset_type"`
	OnBehalfOfContentOwner string `yaml:"on_behalf_of_content_owner" json:"on_behalf_of_content_owner"`
}

type Watermark interface {
	Set()
	Unset()
}

type Option func(*watermark)

func NewWatermark(opts ...Option) Watermark {
	w := &watermark{}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *watermark) Set() {
	file, err := os.Open(w.File)
	if err != nil {
		utils.PrintJSON(w)
		log.Fatalln(errors.Join(errSetWatermark, err))
	}
	defer file.Close()
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

	call := service.Watermarks.Set(w.ChannelId, inVideoBranding).Media(file)
	if w.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.OnBehalfOfContentOwner)
	}

	err = call.Do()
	if err != nil {
		utils.PrintJSON(w)
		log.Fatalln(errors.Join(errSetWatermark, err))
	}

	fmt.Printf("Watermark set for channel %s\n", w.ChannelId)
}

func (w *watermark) Unset() {
	call := service.Watermarks.Unset(w.ChannelId)
	if w.OnBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.OnBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		utils.PrintJSON(w)
		log.Fatalln(errors.Join(errUnsetWatermark, err))
	}

	fmt.Printf("Watermark unset for channel %s\n", w.ChannelId)
}

func WithChannelId(channelId string) Option {
	return func(w *watermark) {
		w.ChannelId = channelId
	}
}

func WithFile(file string) Option {
	return func(w *watermark) {
		w.File = file
	}
}

func WithInVideoPosition(inVideoPosition string) Option {
	return func(w *watermark) {
		w.InVideoPosition = inVideoPosition
	}
}

func WithDurationMs(durationMs uint64) Option {
	return func(w *watermark) {
		w.DurationMs = durationMs
	}
}

func WithOffsetMs(offsetMs uint64) Option {
	return func(w *watermark) {
		w.OffsetMs = offsetMs
	}
}

func WithOffsetType(offsetType string) Option {
	return func(w *watermark) {
		w.OffsetType = offsetType
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(w *watermark) {
		w.OnBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithService(svc *youtube.Service) Option {
	return func(_ *watermark) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential(""),
				auth.WithCacheToken(""),
			)
		}
		service = svc
	}
}

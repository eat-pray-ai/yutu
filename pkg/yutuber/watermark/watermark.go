package watermark

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service           *youtube.Service
	errOpenFile       = errors.New("failed to open file")
	errSetWatermark   = errors.New("failed to set watermark")
	errUnsetWatermark = errors.New("failed to unset watermark")
)

type watermark struct {
	channelId              string
	file                   string
	inVideoPosition        string
	durationMs             uint64
	offsetMs               uint64
	offsetType             string
	onBehalfOfContentOwner string
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

func (w watermark) Set() {
	file, err := os.Open(w.file)
	if err != nil {
		log.Fatalln(errors.Join(errOpenFile, err), w.file)
	}
	defer file.Close()
	inVideoBranding := &youtube.InvideoBranding{
		Position: &youtube.InvideoPosition{},
		Timing:   &youtube.InvideoTiming{},
	}
	if w.inVideoPosition != "" {
		inVideoBranding.Position.Type = "corner"
		inVideoBranding.Position.CornerPosition = w.inVideoPosition
	}
	if w.durationMs != 0 {
		inVideoBranding.Timing.DurationMs = w.durationMs
	}
	if w.offsetMs != 0 {
		inVideoBranding.Timing.OffsetMs = w.offsetMs
	}
	if w.offsetType != "" {
		inVideoBranding.Timing.Type = w.offsetType
	}

	call := service.Watermarks.Set(w.channelId, inVideoBranding).Media(file)
	if w.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.onBehalfOfContentOwner)
	}

	err = call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errSetWatermark, err))
	}
	fmt.Println("Watermark set done")

}

func (w *watermark) Unset() {
	call := service.Watermarks.Unset(w.channelId)
	if w.onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(w.onBehalfOfContentOwner)
	}

	err := call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errUnsetWatermark, err))
	}

	fmt.Println("Watermark unset done")
}

func WithChannelId(channelId string) Option {
	return func(w *watermark) {
		w.channelId = channelId
	}
}

func WithFile(file string) Option {
	return func(w *watermark) {
		w.file = file
	}
}

func WithInVideoPosition(inVideoPosition string) Option {
	return func(w *watermark) {
		w.inVideoPosition = inVideoPosition
	}
}

func WithDurationMs(durationMs uint64) Option {
	return func(w *watermark) {
		w.durationMs = durationMs
	}
}

func WithOffsetMs(offsetMs uint64) Option {
	return func(w *watermark) {
		w.offsetMs = offsetMs
	}
}

func WithOffsetType(offsetType string) Option {
	return func(w *watermark) {
		w.offsetType = offsetType
	}
}

func WithOnBehalfOfContentOwner(onBehalfOfContentOwner string) Option {
	return func(w *watermark) {
		w.onBehalfOfContentOwner = onBehalfOfContentOwner
	}
}

func WithService(svc *youtube.Service) Option {
	return func(w *watermark) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}

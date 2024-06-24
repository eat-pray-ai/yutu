package thumbnail

import (
	"errors"
	"fmt"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service         *youtube.Service
	errSetThumbnail = errors.New("failed to set thumbnail")
)

type thumbnail struct {
	file    string
	videoId string
}

type Thumbnail interface {
	Set()
}

type Option func(*thumbnail)

func NewThumbnail(opts ...Option) Thumbnail {
	t := &thumbnail{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *thumbnail) Set() {
	file, err := os.Open(t.file)
	if err != nil {
		log.Fatalln(errors.Join(errSetThumbnail, err), t.file)
	}
	call := service.Thumbnails.Set(t.videoId).Media(file)
	_, err = call.Do()
	if err != nil {
		log.Fatalln(errors.Join(errSetThumbnail, err))
	}
	fmt.Printf("Thumbnail set for video ID %v\n", t.videoId)
}

func WithVideoId(videoId string) Option {
	return func(t *thumbnail) {
		t.videoId = videoId
	}
}

func WithFile(file string) Option {
	return func(t *thumbnail) {
		t.file = file
	}
}

func WithService(svc *youtube.Service) Option {
	return func(t *thumbnail) {
		if svc != nil {
			service = svc
		} else {
			service = auth.NewY2BService()
		}
	}
}

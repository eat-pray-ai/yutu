package thumbnail

import (
	"errors"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
	"log"
	"os"
)

var (
	service         *youtube.Service
	errSetThumbnail = errors.New("failed to set thumbnail")
)

type thumbnail struct {
	File    string `yaml:"file" json:"file"`
	VideoId string `yaml:"video_id" json:"video_id"`
}

type Thumbnail interface {
	Set(silent bool)
}

type Option func(*thumbnail)

func NewThumbnail(opts ...Option) Thumbnail {
	t := &thumbnail{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *thumbnail) Set(silent bool) {
	file, err := os.Open(t.File)
	if err != nil {
		utils.PrintJSON(t)
		log.Fatalln(errors.Join(errSetThumbnail, err), t.File)
	}
	call := service.Thumbnails.Set(t.VideoId).Media(file)
	res, err := call.Do()
	if err != nil {
		utils.PrintJSON(t)
		log.Fatalln(errors.Join(errSetThumbnail, err))
	}

	if !silent {
		utils.PrintYAML(res)
	}
}

func WithVideoId(videoId string) Option {
	return func(t *thumbnail) {
		t.VideoId = videoId
	}
}

func WithFile(file string) Option {
	return func(t *thumbnail) {
		t.File = file
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

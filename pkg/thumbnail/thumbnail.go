package thumbnail

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
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
	Set(string, string, io.Writer) error
}

type Option func(*thumbnail)

func NewThumbnail(opts ...Option) Thumbnail {
	t := &thumbnail{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *thumbnail) Set(output string, jpath string, writer io.Writer) error {
	file, err := pkg.Root.Open(t.File)
	if err != nil {
		return errors.Join(errSetThumbnail, err)
	}

	call := service.Thumbnails.Set(t.VideoId).Media(file)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errSetThumbnail, err)
	}

	switch output {
	case "json":
		utils.PrintJSON(res, jpath, writer)
	case "yaml":
		utils.PrintYAML(res, jpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Thumbnail set for video %s", t.VideoId)
	}
	return nil
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
	return func(_ *thumbnail) {
		if svc == nil {
			svc = auth.NewY2BService(
				auth.WithCredential("", pkg.Root.FS()),
				auth.WithCacheToken("", pkg.Root.FS()),
			).GetService()
		}
		service = svc
	}
}

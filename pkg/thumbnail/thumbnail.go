// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

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
	errSetThumbnail = errors.New("failed to set thumbnail")
)

type Thumbnail struct {
	File     string `yaml:"file" json:"file"`
	VideoId  string `yaml:"video_id" json:"video_id"`
	Output   string `yaml:"output" json:"output"`
	Jsonpath string `yaml:"jsonpath" json:"jsonpath"`

	service *youtube.Service
}

type IThumbnail interface {
	Set(io.Writer) error
	preRun()
}

type Option func(*Thumbnail)

func NewThumbnail(opts ...Option) IThumbnail {
	t := &Thumbnail{}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Thumbnail) preRun() {
	if t.service == nil {
		t.service = auth.NewY2BService(
			auth.WithCredential("", pkg.Root.FS()),
			auth.WithCacheToken("", pkg.Root.FS()),
		).GetService()
	}
}

func (t *Thumbnail) Set(writer io.Writer) error {
	t.preRun()
	file, err := pkg.Root.Open(t.File)
	if err != nil {
		return errors.Join(errSetThumbnail, err)
	}

	call := t.service.Thumbnails.Set(t.VideoId).Media(file)
	res, err := call.Do()
	if err != nil {
		return errors.Join(errSetThumbnail, err)
	}

	switch t.Output {
	case "json":
		utils.PrintJSON(res, t.Jsonpath, writer)
	case "yaml":
		utils.PrintYAML(res, t.Jsonpath, writer)
	case "silent":
	default:
		_, _ = fmt.Fprintf(writer, "Thumbnail set for video %s", t.VideoId)
	}
	return nil
}

func WithVideoId(videoId string) Option {
	return func(t *Thumbnail) {
		t.VideoId = videoId
	}
}

func WithFile(file string) Option {
	return func(t *Thumbnail) {
		t.File = file
	}
}

func WithOutput(output string) Option {
	return func(t *Thumbnail) {
		t.Output = output
	}
}

func WithJsonpath(jsonpath string) Option {
	return func(t *Thumbnail) {
		t.Jsonpath = jsonpath
	}
}

func WithService(svc *youtube.Service) Option {
	return func(t *Thumbnail) {
		t.service = svc
	}
}

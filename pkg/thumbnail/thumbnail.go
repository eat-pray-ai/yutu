// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"errors"
	"fmt"
	"io"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/eat-pray-ai/yutu/pkg/utils"
)

var (
	errSetThumbnail = errors.New("failed to set thumbnail")
)

type Thumbnail struct {
	*common.Fields
	File    string `yaml:"file" json:"file,omitempty"`
	VideoId string `yaml:"video_id" json:"video_id,omitempty"`
}

type IThumbnail interface {
	Set(io.Writer) error
}

type Option func(*Thumbnail)

func NewThumbnail(opts ...Option) IThumbnail {
	t := &Thumbnail{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *Thumbnail) Set(writer io.Writer) error {
	t.EnsureService()
	file, err := pkg.Root.Open(t.File)
	if err != nil {
		return errors.Join(errSetThumbnail, err)
	}

	call := t.Service.Thumbnails.Set(t.VideoId).Media(file)
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

var (
	WithOutput   = common.WithOutput[*Thumbnail]
	WithJsonpath = common.WithJsonpath[*Thumbnail]
	WithService  = common.WithService[*Thumbnail]
)

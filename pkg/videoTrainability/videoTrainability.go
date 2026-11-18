// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoTrainability

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideoTrainability = errors.New("failed to get video trainability")
)

type VideoTrainability struct {
	common.Fields
	VideoId string
}

type IVideoTrainability[T any] interface {
	Get() (*T, error)
	List(io.Writer) error
}

type Option func(*VideoTrainability)

func WithVideoId(videoId string) Option {
	return func(v *VideoTrainability) {
		v.VideoId = videoId
	}
}

func WithOutput(output string) Option {
	return func(v *VideoTrainability) {
		v.Output = output
	}
}

func WithService(svc *youtube.Service) Option {
	return func(v *VideoTrainability) {
		v.Service = svc
	}
}

func NewVideoTrainability(opts ...Option) IVideoTrainability[youtube.VideoTrainability] {
	v := &VideoTrainability{Fields: common.Fields{}}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

func (v *VideoTrainability) EnsureService() error {
	if v.Service != nil {
		return nil
	}

	ctx := v.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	svc, err := youtube.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	v.Service = svc
	return nil
}

func (v *VideoTrainability) Get() (*youtube.VideoTrainability, error) {
	if err := v.EnsureService(); err != nil {
		return nil, err
	}

	call := v.Service.VideoTrainability.Get().Id(v.VideoId)
	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoTrainability, err)
	}
	return res, nil
}

func (v *VideoTrainability) List(writer io.Writer) error {
	res, err := v.Get()
	if err != nil {
		return err
	}

	return common.PrintList(
		v.Output,
		[]*youtube.VideoTrainability{res},
		writer,
		table.Row{"Video ID", "Permitted"},
		func(vt *youtube.VideoTrainability) table.Row {
			return table.Row{vt.VideoId, strings.Join(vt.Permitted, ", ")}
		},
	)
}

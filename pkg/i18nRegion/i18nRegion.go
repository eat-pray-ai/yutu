// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetI18nRegion = errors.New("failed to get i18n region")
)

type I18nRegion struct {
	*common.Fields
}

type II18nRegion[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*I18nRegion)

func NewI18nRegion(opts ...Option) II18nRegion[youtube.I18nRegion] {
	i := &I18nRegion{Fields: &common.Fields{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *I18nRegion) Get() ([]*youtube.I18nRegion, error) {
	if err := i.EnsureService(); err != nil {
		return nil, err
	}
	call := i.Service.I18nRegions.List(i.Parts)
	if i.Hl != "" {
		call = call.Hl(i.Hl)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetI18nRegion, err)
	}

	return res.Items, nil
}

func (i *I18nRegion) List(writer io.Writer) error {
	regions, err := i.Get()
	if err != nil {
		return err
	}

	common.PrintList(i.Output, regions, writer, table.Row{"ID", "Gl", "Name"}, func(r *youtube.I18nRegion) table.Row {
		return table.Row{r.Id, r.Snippet.Gl, r.Snippet.Name}
	})
	return nil
}

var (
	WithHl = common.WithHl[*I18nRegion]
	WithParts   = common.WithParts[*I18nRegion]
	WithOutput  = common.WithOutput[*I18nRegion]
	WithService = common.WithService[*I18nRegion]
)

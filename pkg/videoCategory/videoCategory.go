// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"errors"
	"io"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

var (
	errGetVideoCategory = errors.New("failed to get video categoryId")
)

type VideoCategory struct {
	*common.Fields
	RegionCode string `yaml:"region_code" json:"region_code,omitempty"`
}

type IVideoCategory[T any] interface {
	Get() ([]*T, error)
	List(io.Writer) error
}

type Option func(*VideoCategory)

func NewVideoCategory(opt ...Option) IVideoCategory[youtube.VideoCategory] {
	vc := &VideoCategory{Fields: &common.Fields{}}
	for _, o := range opt {
		o(vc)
	}
	return vc
}

func (vc *VideoCategory) Get() ([]*youtube.VideoCategory, error) {
	if err := vc.EnsureService(); err != nil {
		return nil, err
	}
	call := vc.Service.VideoCategories.List(vc.Parts)
	if len(vc.Ids) > 0 {
		call = call.Id(vc.Ids...)
	}
	if vc.Hl != "" {
		call = call.Hl(vc.Hl)
	}
	if vc.RegionCode != "" {
		call = call.RegionCode(vc.RegionCode)
	}

	res, err := call.Do()
	if err != nil {
		return nil, errors.Join(errGetVideoCategory, err)
	}

	return res.Items, nil
}

func (vc *VideoCategory) List(writer io.Writer) error {
	categories, err := vc.Get()
	if err != nil {
		return err
	}

	common.PrintList(vc.Output, categories, writer, table.Row{"ID", "Title", "Assignable"}, func(c *youtube.VideoCategory) table.Row {
		return table.Row{c.Id, c.Snippet.Title, c.Snippet.Assignable}
	})
	return nil
}

func WithRegionCode(regionCode string) Option {
	return func(vc *VideoCategory) {
		vc.RegionCode = regionCode
	}
}

var (
	WithIds     = common.WithIds[*VideoCategory]
	WithHl      = common.WithHl[*VideoCategory]
	WithParts   = common.WithParts[*VideoCategory]
	WithOutput  = common.WithOutput[*VideoCategory]
	WithService = common.WithService[*VideoCategory]
)
